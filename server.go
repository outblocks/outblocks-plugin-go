package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/outblocks/outblocks-plugin-go/env"
	"github.com/outblocks/outblocks-plugin-go/log"
)

const ProtocolV1 = "v1"

type Server struct {
	quit chan struct{}
	log  log.Logger
	env  env.Enver
	wg   sync.WaitGroup
}

func NewServer() *Server {
	rand.Seed(time.Now().UnixNano())

	return &Server{
		quit: make(chan struct{}),
		log:  log.NewLogger(),
		env:  env.NewEnv(),
	}
}

func (s *Server) handleConnection(ctx context.Context, handler *ReqHandler, c net.Conn) {
	err := handler.Handle(ctx, s.log, c)
	_ = c.Close()

	if err != nil {
		s.log.Fatalln(err)
	}

	s.wg.Done()
}

func (s *Server) Start(handler *ReqHandler) error {
	handshake := Handshake{
		Protocol: ProtocolV1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := net.Listen("tcp4", "")
	if err != nil {
		s.log.Fatalln(err)
	}

	defer l.Close()

	handshake.Addr = l.Addr().String()

	out, err := json.Marshal(handshake)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

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

		close(s.quit)
		l.Close()
	}()

	for {
		c, err := l.Accept()

		if err != nil {
			select {
			case <-s.quit:
				cancel()

				s.wg.Wait()

				if handler.Cleanup != nil {
					return handler.Cleanup()
				}

				return nil
			default:
				s.log.Fatalln(err)
			}
		}

		s.wg.Add(1)

		go s.handleConnection(ctx, handler, c)
	}
}

func (s *Server) Log() log.Logger {
	return s.log
}

func (s *Server) Env() env.Enver {
	return s.env
}
