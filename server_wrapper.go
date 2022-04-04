package plugin

import (
	"context"

	"github.com/outblocks/outblocks-plugin-go/env"
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/log"
	"github.com/outblocks/outblocks-plugin-go/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type DNSPluginHandler interface {
	GetDomainInfo(ctx context.Context, in *apiv1.DomainInfoRequest) (*apiv1.DomainInfoResponse, error)
	PlanDNS(context.Context, *registry.Registry, *apiv1.PlanDNSRequest) (*apiv1.PlanDNSResponse, error)
	ApplyDNS(*apiv1.ApplyDNSRequest, *registry.Registry, apiv1.DNSPluginService_ApplyDNSServer) error
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
	env  env.Enver
	conn *grpc.ClientConn
	BasicPluginHandler
}

func (s *basicPluginHandlerWrapper) Init(ctx context.Context, req *apiv1.InitRequest) (*apiv1.InitResponse, error) {
	conn, err := grpc.DialContext(ctx, req.HostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	s.conn = conn
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

type dnsPluginHandlerWrapper struct {
	DNSPluginHandler

	RegistryOptions RegistryOptions
}

func (s *dnsPluginHandlerWrapper) createRegistry(read bool) *registry.Registry {
	reg := registry.NewRegistry(&registry.Options{
		Read:            read,
		AllowDuplicates: s.RegistryOptions.AllowDuplicates,
	})

	return reg
}
func (s *dnsPluginHandlerWrapper) PlanDNS(ctx context.Context, r *apiv1.PlanDNSRequest) (*apiv1.PlanDNSResponse, error) {
	reg := s.createRegistry(r.Verify)
	return s.DNSPluginHandler.PlanDNS(ctx, reg, r)
}

func (s *dnsPluginHandlerWrapper) ApplyDNS(r *apiv1.ApplyDNSRequest, stream apiv1.DNSPluginService_ApplyDNSServer) error {
	reg := s.createRegistry(false)
	return s.DNSPluginHandler.ApplyDNS(r, reg, stream)
}

func DefaultRegistryApplyCallback(stream apiv1.DeployPluginService_ApplyServer) func(*apiv1.ApplyAction) {
	return func(a *apiv1.ApplyAction) {
		_ = stream.Send(&apiv1.ApplyResponse{
			Response: &apiv1.ApplyResponse_Action{
				Action: &apiv1.ApplyActionResponse{
					Actions: []*apiv1.ApplyAction{a},
				},
			},
		})
	}
}
