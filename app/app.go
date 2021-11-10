package app

import (
	"context"
	"fmt"

	"github.com/tkeel-io/kit/transport/grpc"
	"github.com/tkeel-io/kit/transport/http"
)

type App struct {
	Name       string
	httpServer *http.Server
	grpcServer *grpc.Server
}

func New(name, httpAddr, grpcAddr string) *App {
	return &App{
		Name:       name,
		httpServer: http.NewServer(httpAddr),
		grpcServer: grpc.NewServer(grpcAddr),
	}
}

func (a *App) Run(ctx context.Context) error {
	if err := a.grpcServer.Start(ctx); err != nil {
		return fmt.Errorf("error start grpc: %w", err)
	}
	if err := a.httpServer.Start(ctx); err != nil {
		return fmt.Errorf("error start http: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.grpcServer.Stop(ctx); err != nil {
		return fmt.Errorf("error stop grpc: %w", err)
	}
	if err := a.httpServer.Stop(ctx); err != nil {
		return fmt.Errorf("error Stop http: %w", err)
	}
	return nil
}
