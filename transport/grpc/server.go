package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/tkeel-io/kit/log"
	"github.com/tkeel-io/kit/transport"
	"google.golang.org/grpc"
)

const DefaultPort = ":31233"

type Server struct {
	Addr string
	srv  *grpc.Server
}

func NewServer(addr string) *Server {
	if addr == "" {
		addr = DefaultPort
	}
	return &Server{
		Addr: addr,
		srv:  grpc.NewServer(),
	}
}

func (s *Server) GetServe() *grpc.Server {
	return s.srv
}

func (s *Server) Type() transport.Type {
	return transport.TypeGRPC
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("error listen addr: %w", err)
	}
	log.Debugf("GRPC Server listen: %s", s.Addr)
	go func() {
		if err := s.srv.Serve(l); err != nil {
			log.Errorf("error http serve: %w", err)
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.srv.Stop()
	return nil
}
