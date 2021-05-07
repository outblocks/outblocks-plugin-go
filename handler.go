package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/outblocks/outblocks-plugin-go/log"
)

type ReqHandler struct {
	Init               func(ctx context.Context, r *InitRequest) (Response, error)
	InitInteractive    func(ctx context.Context, r *InitRequest, in chan<- Request, out <-chan Response) error
	Start              func(ctx context.Context, r *StartRequest) (Response, error)
	StartInteractive   func(ctx context.Context, r *StartRequest, in chan<- Request, out <-chan Response) error
	Plan               func(ctx context.Context, r *PlanRequest) (Response, error)
	PlanInteractive    func(ctx context.Context, r *PlanRequest, in chan<- Request, out <-chan Response) error
	Apply              func(ctx context.Context, r *ApplyRequest) (Response, error)
	AppleInteractive   func(ctx context.Context, r *ApplyRequest, in chan<- Request, out <-chan Response) error
	Run                func(ctx context.Context, r *RunRequest) (Response, error)
	RunInteractive     func(ctx context.Context, r *RunRequest, in chan<- Request, out <-chan Response) error
	Command            func(ctx context.Context, r *CommandRequest) (Response, error)
	CommandInteractive func(ctx context.Context, r *CommandRequest, in chan<- Request, out <-chan Response) error

	// Cleanup
	Cleanup func() error

	// State handlers.
	GetState  func(ctx context.Context, r *GetStateRequest) (*GetStateResponse, error)
	SaveState func(ctx context.Context, r *SaveStateRequest) (*SaveStateResponse, error)
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

	return
}

func (h *ReqHandler) handleInteractive(ctx context.Context, logger log.Logger, c net.Conn, r *bufio.Reader, req Request) (bool, error) {
	var handler func() error

	in := make(chan Request)
	out := make(chan Response)

	switch v := req.(type) {
	case *InitRequest:
		if h.InitInteractive != nil {
			handler = func() error { return h.InitInteractive(ctx, v, in, out) }
		}
	case *StartRequest:
		if h.StartInteractive != nil {
			handler = func() error { return h.StartInteractive(ctx, v, in, out) }
		}
	case *PlanRequest:
		if h.PlanInteractive != nil {
			handler = func() error { return h.PlanInteractive(ctx, v, in, out) }
		}
	case *ApplyRequest:
		if h.AppleInteractive != nil {
			handler = func() error { return h.AppleInteractive(ctx, v, in, out) }
		}
	case *RunRequest:
		if h.RunInteractive != nil {
			handler = func() error { return h.RunInteractive(ctx, v, in, out) }
		}
	case *CommandRequest:
		if h.CommandInteractive != nil {
			handler = func() error { return h.CommandInteractive(ctx, v, in, out) }
		}
	}

	if handler == nil {
		return false, nil
	}

	var (
		err error
		wg  sync.WaitGroup
	)

	errCh := make(chan error)
	wg.Add(2)

	go func() {
		for {
			req, err := readRequest(logger, r)
			if err != nil {
				errCh <- err
				close(in)
				wg.Done()
				return
			}

			in <- req
		}
	}()

	go func() {
		for {
			select {
			case res, ok := <-out:
				if !ok {
					c.Close()
					return
				}
				err := writeResponse(c, res)
				if err != nil {
					errCh <- err
					break
				}
			}
		}

		wg.Done()
	}()

	if err := handler(); err != nil {
		return true, err
	}

	wg.Wait()

	select {
	case err = <-errCh:
	}

	return true, err
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
		return err
	}

	if res != nil {
		return writeResponse(c, res)
	}

	// Handle interactive responses.
	handled, err := h.handleInteractive(ctx, logger, c, r, req)
	if err != nil {
		return err
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

	if sr, ok := res.(StreamingResponse); ok {
		header.IsStreaming = true
		header.IsFinal = sr.IsFinal()
	}

	// Send header.
	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	if _, err := w.Write([]byte{'\n'}); err != nil {
		return err
	}

	// Send response itself.
	data, err = json.Marshal(res)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})

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
