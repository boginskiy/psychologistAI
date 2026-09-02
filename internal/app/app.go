package app

import (
	"context"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/logger"
	"github.com/boginskiy/psychologistAI/internal/router"
	"github.com/boginskiy/psychologistAI/internal/server"
	"github.com/boginskiy/psychologistAI/internal/service"
)

type App struct {
	Cfg  config.Config
	Logg logger.Logger

	Server server.Server
	Router router.Router
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initModules(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.Server.Run(ctx, a.Router.Run())
}

func (a *App) initModules(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initHandlers,
		a.initRouter,
		a.initServer,
	}

	for _, init := range inits {
		err := init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initHandlers(ctx context.Context) error {
	// Infra services
	validater := service.NewValidService(ctx)
	response := response.NewResp()

	// Services
	userService := service.NewUserServ(ctx, validater)

	// Handlers
	userHandler := handlers.NewUserHandler("/user", userService, response)

	// Router
	a.Router.RegisterRoutes(userHandler)
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	a.Router = router.NewRouterChi(ctx, "/api/v1")
	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	// NEED create
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	// NEED create
	return nil
}

func (a *App) initServer(ctx context.Context) error {
	a.Server = server.NewServerHTTP(ctx, server.ConfigServer)
	return nil
}
