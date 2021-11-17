package http

import (
	"context"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/tkeel-io/kit/log"
	"github.com/tkeel-io/kit/transport"
)

const (
	ImportAndUsed = true
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
	restful.RegisterEntityAccessor("application/x-www-form-urlencoded", FormEntityReadWriter{})
	c.EnableContentEncoding(true)
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
	log.Debugf("HTTP Server listen: %s", s.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Errorf("error http serve: %w", err)
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
