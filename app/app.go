package app

import (
	"context"
	"fmt"

	"github.com/tkeel-io/kit/log"
	"github.com/tkeel-io/kit/transport"
)

type App struct {
	Name       string
	serverList []transport.Server
}

func New(name string, conf *log.Conf, srv ...transport.Server) *App {
	if err := log.InitLoggerByConf(conf); err != nil {
		panic(err)
	}
	app := &App{
		Name:       name,
		serverList: srv,
	}
	return app
}

func (a *App) Run(ctx context.Context) error {
	for _, v := range a.serverList {
		if err := v.Start(ctx); err != nil {
			return fmt.Errorf("error start server(%s): %w", v.Type(), err)
		}
	}
	log.Infof("app %s running", a.Name)
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	for _, v := range a.serverList {
		if err := v.Stop(ctx); err != nil {
			return fmt.Errorf("error stop server(%s): %w", v.Type(), err)
		}
	}
	return nil
}
