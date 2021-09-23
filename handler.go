package plugin

import (
	"context"
	"net"

	"github.com/outblocks/outblocks-plugin-go/log"
)

type ReqHandler struct {
	ProjectInit            func(ctx context.Context, r *ProjectInitRequest) (Response, error)
	ProjectInitInteractive func(ctx context.Context, r *ProjectInitRequest, stream *ReceiverStream) error
	Start                  func(ctx context.Context, r *StartRequest) (Response, error)
	StartInteractive       func(ctx context.Context, r *StartRequest, stream *ReceiverStream) error
	Plan                   func(ctx context.Context, r *PlanRequest) (Response, error)
	PlanInteractive        func(ctx context.Context, r *PlanRequest, stream *ReceiverStream) error
	Apply                  func(ctx context.Context, r *ApplyRequest) (Response, error)
	ApplyInteractive       func(ctx context.Context, r *ApplyRequest, stream *ReceiverStream) error
	Run                    func(ctx context.Context, r *RunRequest) (Response, error)
	RunInteractive         func(ctx context.Context, r *RunRequest, stream *ReceiverStream) error
	Command                func(ctx context.Context, r *CommandRequest) (Response, error)
	CommandInteractive     func(ctx context.Context, r *CommandRequest, stream *ReceiverStream) error

	// Cleanup
	Cleanup func() error

	// State handlers.
	GetState    func(ctx context.Context, r *GetStateRequest) (Response, error)
	SaveState   func(ctx context.Context, r *SaveStateRequest) (Response, error)
	ReleaseLock func(ctx context.Context, r *ReleaseLockRequest) (Response, error)
}

func (h *ReqHandler) handleSync(ctx context.Context, req Request) (res Response, err error) {
	switch v := req.(type) {
	case *ProjectInitRequest:
		if h.ProjectInit != nil {
			res, err = h.ProjectInit(ctx, v)
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

func (h *ReqHandler) interactiveHandler(ctx context.Context, req Request, stream *ReceiverStream) func() error {
	switch v := req.(type) {
	case *ProjectInitRequest:
		if h.ProjectInitInteractive != nil {
			return func() error { return h.ProjectInitInteractive(ctx, v, stream) }
		}
	case *StartRequest:
		if h.StartInteractive != nil {
			return func() error { return h.StartInteractive(ctx, v, stream) }
		}
	case *PlanRequest:
		if h.PlanInteractive != nil {
			return func() error { return h.PlanInteractive(ctx, v, stream) }
		}
	case *ApplyRequest:
		if h.ApplyInteractive != nil {
			return func() error { return h.ApplyInteractive(ctx, v, stream) }
		}
	case *RunRequest:
		if h.RunInteractive != nil {
			return func() error { return h.RunInteractive(ctx, v, stream) }
		}
	case *CommandRequest:
		if h.CommandInteractive != nil {
			return func() error { return h.CommandInteractive(ctx, v, stream) }
		}
	}

	return nil
}

func (h *ReqHandler) handleInteractive(ctx context.Context, req Request, stream *ReceiverStream) (handled bool, err error) {
	handler := h.interactiveHandler(ctx, req, stream)
	if handler == nil {
		return false, nil
	}

	err = handler()
	if err != nil {
		return true, err
	}

	return true, nil
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
	stream := NewReceiverStream(c)
	defer stream.Close()

	req, err := stream.Recv()
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
		return stream.Send(res)
	}

	// Handle interactive responses.
	handled, err := h.handleInteractive(ctx, req, stream)
	if err != nil {
		handleError(c, err)

		return nil
	}

	if handled {
		return nil
	}

	return stream.Send(&UnhandledResponse{})
}
