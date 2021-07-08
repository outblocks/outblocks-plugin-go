package plugin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/outblocks/outblocks-plugin-go/log"
)

type ReqHandler struct {
	Init               func(ctx context.Context, r *InitRequest) (Response, error)
	InitInteractive    func(ctx context.Context, r *InitRequest, in <-chan Request, out chan<- Response) error
	Start              func(ctx context.Context, r *StartRequest) (Response, error)
	StartInteractive   func(ctx context.Context, r *StartRequest, in <-chan Request, out chan<- Response) error
	Plan               func(ctx context.Context, r *PlanRequest) (Response, error)
	PlanInteractive    func(ctx context.Context, r *PlanRequest, in <-chan Request, out chan<- Response) error
	Apply              func(ctx context.Context, r *ApplyRequest) (Response, error)
	ApplyInteractive   func(ctx context.Context, r *ApplyRequest, in <-chan Request, out chan<- Response) error
	Run                func(ctx context.Context, r *RunRequest) (Response, error)
	RunInteractive     func(ctx context.Context, r *RunRequest, in <-chan Request, out chan<- Response) error
	Command            func(ctx context.Context, r *CommandRequest) (Response, error)
	CommandInteractive func(ctx context.Context, r *CommandRequest, in <-chan Request, out chan<- Response) error

	// Cleanup
	Cleanup func() error

	// State handlers.
	GetState    func(ctx context.Context, r *GetStateRequest) (Response, error)
	SaveState   func(ctx context.Context, r *SaveStateRequest) (Response, error)
	ReleaseLock func(ctx context.Context, r *ReleaseLockRequest) (Response, error)
}

func (h *ReqHandler) handleSync(ctx context.Context, req Request) (res Response, err error) {
	switch v := req.(type) {
	case *InitRequest:
		if h.Init != nil {
			res, err = h.Init(ctx, v)
		}
	case *StartRequest:
		if h.Start != nil {
			res, err = h.Start(ctx, v)
		}
	case *GetStateRequest:
		if h.GetState != nil {
			res, err = h.GetState(ctx, v)
		}
	case *SaveStateRequest:
		if h.SaveState != nil {
			res, err = h.SaveState(ctx, v)
		}
	case *ReleaseLockRequest:
		if h.ReleaseLock != nil {
			res, err = h.ReleaseLock(ctx, v)
		}
	case *PlanRequest:
		if h.Plan != nil {
			res, err = h.Plan(ctx, v)
		}
	case *ApplyRequest:
		if h.Apply != nil {
			res, err = h.Apply(ctx, v)
		}
	case *RunRequest:
		if h.Run != nil {
			res, err = h.Run(ctx, v)
		}
	case *CommandRequest:
		if h.Command != nil {
			res, err = h.Command(ctx, v)
		}
	}

	return res, err
}

func (h *ReqHandler) interactiveHandler(ctx context.Context, req Request, in chan Request, out chan Response) func() error {
	switch v := req.(type) {
	case *InitRequest:
		if h.InitInteractive != nil {
			return func() error { return h.InitInteractive(ctx, v, in, out) }
		}
	case *StartRequest:
		if h.StartInteractive != nil {
			return func() error { return h.StartInteractive(ctx, v, in, out) }
		}
	case *PlanRequest:
		if h.PlanInteractive != nil {
			return func() error { return h.PlanInteractive(ctx, v, in, out) }
		}
	case *ApplyRequest:
		if h.ApplyInteractive != nil {
			return func() error { return h.ApplyInteractive(ctx, v, in, out) }
		}
	case *RunRequest:
		if h.RunInteractive != nil {
			return func() error { return h.RunInteractive(ctx, v, in, out) }
		}
	case *CommandRequest:
		if h.CommandInteractive != nil {
			return func() error { return h.CommandInteractive(ctx, v, in, out) }
		}
	}

	return nil
}

func (h *ReqHandler) handleInteractive(ctx context.Context, logger log.Logger, c net.Conn, r *bufio.Reader, req Request) (handled bool, err error) {
	var handler func() error

	in := make(chan Request)
	out := make(chan Response)

	handler = h.interactiveHandler(ctx, req, in, out)
	if handler == nil {
		return false, nil
	}

	writeWait := make(chan struct{})

	errCh := make(chan error, 2)

	go func() {
		for {
			req, err := readRequest(logger, r)
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}

				close(in)

				return
			}

			in <- req
		}
	}()

	go func() {
		defer close(writeWait)

		for res := range out {
			err := writeResponse(c, res)
			if err != nil {
				errCh <- err

				return
			}
		}
	}()

	err = handler()

	close(out)
	<-writeWait

	if err != nil {
		return true, err
	}

	select {
	case err = <-errCh:
	default:
	}

	return true, err
}

func handleError(c net.Conn, err error) {
	if r, ok := err.(Response); ok {
		_ = writeResponse(c, r)

		return
	}

	_ = writeResponse(c, &ErrorResponse{
		Error: err.Error(),
	})
}

func (h *ReqHandler) Handle(ctx context.Context, logger log.Logger, c net.Conn) error {
	r := bufio.NewReader(c)

	req, err := readRequest(logger, r)
	if err != nil {
		return err
	}

	// Handle sync responses.
	res, err := h.handleSync(ctx, req)
	if err != nil {
		handleError(c, err)

		return nil
	}

	if res != nil {
		return writeResponse(c, res)
	}

	// Handle interactive responses.
	handled, err := h.handleInteractive(ctx, logger, c, r, req)
	if err != nil {
		handleError(c, err)

		return nil
	}

	if handled {
		return nil
	}

	return writeResponse(c, &UnhandledResponse{})
}

func writeResponse(w io.Writer, res Response) error {
	header := &ResponseHeader{
		Type: res.Type(),
	}

	// Prepare header.
	headerData, err := json.Marshal(header)
	if err != nil {
		return err
	}

	// Prepare response itself.
	responseData, err := json.Marshal(res)
	if err != nil {
		return err
	}

	data := headerData
	data = append(data, byte('\n'))
	data = append(data, responseData...)
	data = append(data, byte('\n'))

	_, err = w.Write(data)

	return err
}

func readRequest(logger log.Logger, r *bufio.Reader) (Request, error) {
	data, err := r.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("unable to read header: %w", err)
	}

	var header RequestHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("header unmarshal error: %s", err)
	}

	var req Request

	switch header.Type {
	case RequestTypeInit:
		req = &InitRequest{}
	case RequestTypeStart:
		req = &StartRequest{}
	case RequestTypePlan:
		req = &PlanRequest{}
	case RequestTypeApply:
		req = &ApplyRequest{}
	case RequestTypeRun:
		req = &RunRequest{}
	case RequestTypeGetState:
		req = &GetStateRequest{}
	case RequestTypeSaveState:
		req = &SaveStateRequest{}
	case RequestTypeCommand:
		req = &CommandRequest{}
	case RequestTypeReleaseLock:
		req = &ReleaseLockRequest{}
	case RequestTypePromptAnswer:
		req = &PromptAnswerRequest{}
	default:
		logger.Fatalf("unknown request type: %d\n", header.Type)
	}

	data, err = r.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("unable to read request: %w", err)
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("request unmarshal error: %s", err)
	}

	return req, nil
}
