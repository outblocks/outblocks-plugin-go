package plugin

import (
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
	rand.Seed(time.Now().UnixNano())

	return &Server{
		env: env.NewEnv(),
	}
}

type ServerOptions func(s *Server)

func WithRegistryAllowDuplicates(b bool) ServerOptions {
	return func(s *Server) {
		s.registryOptions.AllowDuplicates = b
	}
}

func (s *Server) serve(handler BasicPluginHandler, opts ...ServerOptions) error {
	for _, opt := range opts {
		opt(s)
	}

	// Disable grpc client logging.
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(io.Discard, io.Discard, io.Discard))

	handshake := Handshake{
		Protocol: ProtocolV1,
	}

	l, err := net.Listen("tcp4", "")
	if err != nil {
		panic(err)
	}

	handshake.Addr = l.Addr().String()

	out, err := json.Marshal(handshake)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	grpcServer := grpc.NewServer()
	apiv1.RegisterBasicPluginServiceServer(grpcServer, &basicPluginHandlerWrapper{BasicPluginHandler: handler, env: s.env})

	if srv, ok := handler.(DeployPluginHandler); ok {
		apiv1.RegisterDeployPluginServiceServer(grpcServer, &deployPluginHandlerWrapper{DeployPluginHandler: srv, RegistryOptions: s.registryOptions})
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

func Serve(handler BasicPluginHandler, opts ...ServerOptions) error {
	return newServer().serve(handler, opts...)
}
