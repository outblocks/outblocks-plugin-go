package plugin

import (
	"context"

	"github.com/outblocks/outblocks-plugin-go/env"
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/log"
	"github.com/outblocks/outblocks-plugin-go/registry"
	"google.golang.org/grpc"
)

type BasicPluginHandler interface {
	Init(context.Context, env.Enver, log.Logger, apiv1.HostServiceClient) error
	Start(context.Context, *apiv1.StartRequest) (*apiv1.StartResponse, error)
	ProjectInit(context.Context, *apiv1.ProjectInitRequest) (*apiv1.ProjectInitResponse, error)
}

type DeployPluginHandler interface {
	Plan(context.Context, *registry.Registry, *apiv1.PlanRequest) (*apiv1.PlanResponse, error)
	Apply(*apiv1.ApplyRequest, *registry.Registry, apiv1.DeployPluginService_ApplyServer) error
}

type CommandPluginHandler interface {
	apiv1.CommandPluginServiceServer
}

type RunPluginHandler interface {
	apiv1.RunPluginServiceServer
}

type StatePluginHandler interface {
	apiv1.StatePluginServiceServer
}

type LockingPluginHandler interface {
	apiv1.LockingPluginServiceServer
}

type Cleanup interface {
	Cleanup() error
}

type basicPluginHandlerWrapper struct {
	env env.Enver
	BasicPluginHandler
}

func (s *basicPluginHandlerWrapper) Init(ctx context.Context, req *apiv1.InitRequest) (*apiv1.InitResponse, error) {
	conn, err := grpc.DialContext(ctx, req.HostAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := apiv1.NewHostServiceClient(conn)
	l := log.NewLogger(cli)

	return &apiv1.InitResponse{}, s.BasicPluginHandler.Init(ctx, s.env, l, cli)
}

type deployPluginHandlerWrapper struct {
	DeployPluginHandler

	RegistryOptions RegistryOptions
}

func (s *deployPluginHandlerWrapper) createRegistry(apps []*apiv1.AppPlan, destroy, read bool) *registry.Registry {
	reg := registry.NewRegistry(&registry.Options{
		Destroy:         destroy,
		Read:            read,
		AllowDuplicates: s.RegistryOptions.AllowDuplicates,
	})

	for _, plan := range apps {
		if plan.Skip {
			reg.SkipAppResources(plan.State.App)
		}
	}

	return reg
}
func (s *deployPluginHandlerWrapper) Plan(ctx context.Context, r *apiv1.PlanRequest) (*apiv1.PlanResponse, error) {
	reg := s.createRegistry(r.Apps, r.Destroy, r.Verify)
	return s.DeployPluginHandler.Plan(ctx, reg, r)
}

func (s *deployPluginHandlerWrapper) Apply(r *apiv1.ApplyRequest, stream apiv1.DeployPluginService_ApplyServer) error {
	reg := s.createRegistry(r.Apps, r.Destroy, false)
	return s.DeployPluginHandler.Apply(r, reg, stream)
}
