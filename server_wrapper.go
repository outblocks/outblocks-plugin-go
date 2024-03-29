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
	PlanDNS(context.Context, *registry.Registry, *apiv1.PlanDNSRequest) (*apiv1.PlanDNSResponse, error)
	ApplyDNS(*apiv1.ApplyDNSRequest, *registry.Registry, apiv1.DNSPluginService_ApplyDNSServer) error
}

type MonitoringPluginHandler interface {
	PlanMonitoring(context.Context, *registry.Registry, *apiv1.PlanMonitoringRequest) (*apiv1.PlanMonitoringResponse, error)
	ApplyMonitoring(*apiv1.ApplyMonitoringRequest, *registry.Registry, apiv1.MonitoringPluginService_ApplyMonitoringServer) error
}

type LogsPluginHandler interface {
	apiv1.LogsPluginServiceServer
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

type DeployHookHandler interface {
	apiv1.DeployHookServiceServer
}

type SecretPluginHandler interface {
	apiv1.SecretPluginServiceServer
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

func createRegistry(opts *registry.Options, apps []*apiv1.AppPlan) *registry.Registry {
	reg := registry.NewRegistry(opts)

	for _, plan := range apps {
		if plan.Skip {
			reg.SkipAppResources(plan.State.App)
		}
	}

	return reg
}

func (s *deployPluginHandlerWrapper) createRegistry(apps []*apiv1.AppPlan, destroy, read bool) *registry.Registry {
	return createRegistry(&registry.Options{
		Read:            read,
		Destroy:         destroy,
		AllowDuplicates: s.RegistryOptions.AllowDuplicates,
	}, apps)
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

func (s *dnsPluginHandlerWrapper) createRegistry(read, destroy bool) *registry.Registry {
	return createRegistry(&registry.Options{
		Read:            read,
		Destroy:         destroy,
		AllowDuplicates: s.RegistryOptions.AllowDuplicates,
	}, nil)
}

func (s *dnsPluginHandlerWrapper) PlanDNS(ctx context.Context, r *apiv1.PlanDNSRequest) (*apiv1.PlanDNSResponse, error) {
	reg := s.createRegistry(r.Verify, r.Destroy)
	return s.DNSPluginHandler.PlanDNS(ctx, reg, r)
}

func (s *dnsPluginHandlerWrapper) ApplyDNS(r *apiv1.ApplyDNSRequest, stream apiv1.DNSPluginService_ApplyDNSServer) error {
	reg := s.createRegistry(false, r.Destroy)
	return s.DNSPluginHandler.ApplyDNS(r, reg, stream)
}

type monitoringPluginHandlerWrapper struct {
	MonitoringPluginHandler

	RegistryOptions RegistryOptions
}

func (s *monitoringPluginHandlerWrapper) createRegistry(read, destroy bool) *registry.Registry {
	return createRegistry(&registry.Options{
		Read:            read,
		Destroy:         destroy,
		AllowDuplicates: s.RegistryOptions.AllowDuplicates,
	}, nil)
}

func (s *monitoringPluginHandlerWrapper) PlanMonitoring(ctx context.Context, r *apiv1.PlanMonitoringRequest) (*apiv1.PlanMonitoringResponse, error) {
	reg := s.createRegistry(r.Verify, r.Destroy)
	return s.MonitoringPluginHandler.PlanMonitoring(ctx, reg, r)
}

func (s *monitoringPluginHandlerWrapper) ApplyMonitoring(r *apiv1.ApplyMonitoringRequest, stream apiv1.MonitoringPluginService_ApplyMonitoringServer) error {
	reg := s.createRegistry(false, r.Destroy)
	return s.MonitoringPluginHandler.ApplyMonitoring(r, reg, stream)
}
