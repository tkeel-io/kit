package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/tkeel-io/kit/transport"
)

const DefaultPort = ":31234"

type Server struct {
	Addr string
	srv  *http.Server

	Container *restful.Container
}

func NewServer(addr string) *Server {
	if addr == "" {
		addr = DefaultPort
	}
	c := restful.NewContainer()
	return &Server{
		Addr:      addr,
		Container: c,
		srv: &http.Server{
			Addr:    addr,
			Handler: c,
		},
	}
}

func (s *Server) Type() transport.Type {
	return transport.TypeHTTP
}

func (s *Server) Start(ctx context.Context) error {
	go func() error {
		if err := s.srv.ListenAndServe(); err != nil {
			return fmt.Errorf("error http serve: %w", err)
		}
		return nil
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
