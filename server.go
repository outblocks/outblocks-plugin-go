package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/outblocks/outblocks-plugin-go/env"
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const ProtocolV1 = "v1"

type RegistryOptions struct {
	AllowDuplicates bool
}

type Server struct {
	env env.Enver

	registryOptions RegistryOptions
}

func newServer() *Server {
	rand.Seed(time.Now().UnixNano()) //nolint

	return &Server{
		env: env.NewEnv(),
	}
}

type ServerOption func(*Server)

func WithRegistryAllowDuplicates(b bool) ServerOption {
	return func(s *Server) {
		s.registryOptions.AllowDuplicates = b
	}
}

func (s *Server) serve(handler BasicPluginHandler, opts ...ServerOption) error {
	for _, opt := range opts {
		opt(s)
	}

	// Disable grpc client logging.
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(io.Discard, io.Discard, io.Discard))

	handshake := Handshake{
		Protocol: ProtocolV1,
	}

	lCfg := net.ListenConfig{}

	l, err := lCfg.Listen(context.TODO(), "tcp4", "")
	if err != nil {
		panic(err)
	}

	handshake.Addr = l.Addr().String()

	out, err := json.Marshal(handshake)
	if err != nil {
		return err
	}

	fmt.Println(string(out)) //nolint:forbidigo

	grpcServer := grpc.NewServer()
	basicWrapper := &basicPluginHandlerWrapper{BasicPluginHandler: handler, env: s.env}
	apiv1.RegisterBasicPluginServiceServer(grpcServer, basicWrapper)

	if srv, ok := handler.(DeployPluginHandler); ok {
		apiv1.RegisterDeployPluginServiceServer(grpcServer, &deployPluginHandlerWrapper{DeployPluginHandler: srv, RegistryOptions: s.registryOptions})
	}

	if srv, ok := handler.(DNSPluginHandler); ok {
		apiv1.RegisterDNSPluginServiceServer(grpcServer, &dnsPluginHandlerWrapper{DNSPluginHandler: srv, RegistryOptions: s.registryOptions})
	}

	if srv, ok := handler.(MonitoringPluginHandler); ok {
		apiv1.RegisterMonitoringPluginServiceServer(grpcServer, &monitoringPluginHandlerWrapper{MonitoringPluginHandler: srv, RegistryOptions: s.registryOptions})
	}

	if srv, ok := handler.(LogsPluginHandler); ok {
		apiv1.RegisterLogsPluginServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(CommandPluginHandler); ok {
		apiv1.RegisterCommandPluginServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(RunPluginHandler); ok {
		apiv1.RegisterRunPluginServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(StatePluginHandler); ok {
		apiv1.RegisterStatePluginServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(LockingPluginHandler); ok {
		apiv1.RegisterLockingPluginServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(DeployHookHandler); ok {
		apiv1.RegisterDeployHookServiceServer(grpcServer, srv)
	}

	if srv, ok := handler.(SecretPluginHandler); ok {
		apiv1.RegisterSecretPluginServiceServer(grpcServer, srv)
	}

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-ch

			if sig == syscall.SIGTERM {
				break
			}
		}

		grpcServer.GracefulStop()

		if basicWrapper.conn != nil {
			basicWrapper.conn.Close()
		}
	}()

	err = grpcServer.Serve(l)

	if pc, ok := handler.(Cleanup); ok {
		errCleanup := pc.Cleanup()
		if errCleanup != nil {
			return errCleanup
		}
	}

	return err
}

func Serve(handler BasicPluginHandler, opts ...ServerOption) error {
	return newServer().serve(handler, opts...)
}
