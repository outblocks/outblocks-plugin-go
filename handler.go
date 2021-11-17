package plugin

import (
	"context"
	"net"

	"github.com/outblocks/outblocks-plugin-go/log"
	"github.com/outblocks/outblocks-plugin-go/registry"
)

type ReqHandlerOptions struct {
	RegistryAllowDuplicates bool
}

type ReqHandler struct {
	ProjectInit            func(ctx context.Context, r *ProjectInitRequest) (Response, error)
	ProjectInitInteractive func(ctx context.Context, r *ProjectInitRequest, stream *ReceiverStream) error
	Start                  func(ctx context.Context, r *StartRequest) (Response, error)
	StartInteractive       func(ctx context.Context, r *StartRequest, stream *ReceiverStream) error
	Plan                   func(ctx context.Context, r *PlanRequest, reg *registry.Registry) (Response, error)
	PlanInteractive        func(ctx context.Context, r *PlanRequest, reg *registry.Registry, stream *ReceiverStream) error
	Apply                  func(ctx context.Context, r *ApplyRequest, reg *registry.Registry) (Response, error)
	ApplyInteractive       func(ctx context.Context, r *ApplyRequest, reg *registry.Registry, stream *ReceiverStream) error
	Run                    func(ctx context.Context, r *RunRequest) (Response, error)
	RunInteractive         func(ctx context.Context, r *RunRequest, stream *ReceiverStream) error
	Command                func(ctx context.Context, r *CommandRequest) (Response, error)
	CommandInteractive     func(ctx context.Context, r *CommandRequest, stream *ReceiverStream) error

	// Cleanup
	Cleanup func() error

	// State handlers.
	GetState         func(ctx context.Context, r *GetStateRequest) (Response, error)
	SaveState        func(ctx context.Context, r *SaveStateRequest) (Response, error)
	ReleaseStateLock func(ctx context.Context, r *ReleaseStateLockRequest) (Response, error)

	// Locking.
	AcquireLocks func(ctx context.Context, r *AcquireLocksRequest) (Response, error)
	ReleaseLocks func(ctx context.Context, r *ReleaseLocksRequest) (Response, error)

	Options ReqHandlerOptions
}

func (h *ReqHandler) handleSync(ctx context.Context, req Request) (res Response, err error) { // nolint: gocyclo
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
	case *ReleaseStateLockRequest:
		if h.ReleaseStateLock != nil {
			res, err = h.ReleaseStateLock(ctx, v)
		}
	case *AcquireLocksRequest:
		if h.AcquireLocks != nil {
			res, err = h.AcquireLocks(ctx, v)
		}
	case *ReleaseLocksRequest:
		if h.ReleaseLocks != nil {
			res, err = h.ReleaseLocks(ctx, v)
		}
	case *PlanRequest:
		if h.Plan != nil {
			res, err = h.Plan(ctx, v, h.CreateRegistry(&v.DeployBaseRequest, v.Verify))
		}
	case *ApplyRequest:
		if h.Apply != nil {
			res, err = h.Apply(ctx, v, h.CreateRegistry(&v.DeployBaseRequest, false))
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
			return func() error {
				return h.PlanInteractive(ctx, v, h.CreateRegistry(&v.DeployBaseRequest, v.Verify), stream)
			}
		}
	case *ApplyRequest:
		if h.ApplyInteractive != nil {
			return func() error { return h.ApplyInteractive(ctx, v, h.CreateRegistry(&v.DeployBaseRequest, false), stream) }
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

func (h *ReqHandler) CreateRegistry(v *DeployBaseRequest, read bool) *registry.Registry {
	reg := registry.NewRegistry(&registry.Options{
		Destroy:         v.Destroy,
		Read:            read,
		AllowDuplicates: h.Options.RegistryAllowDuplicates,
	})

	for _, app := range v.Apps {
		if app.Skip {
			reg.SkipAppResources(&app.App.App)
		}
	}

	return reg
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
