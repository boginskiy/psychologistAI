package app

import (
	"context"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/logger"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/repository/userrepo"
	"github.com/boginskiy/psychologistAI/internal/router"
	"github.com/boginskiy/psychologistAI/internal/server"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/internal/service/infra"
	"github.com/boginskiy/psychologistAI/pkg/cookie"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
	"github.com/boginskiy/psychologistAI/pkg/security"
)

const PathCountryDB = "GeoLite2-Country.mmdb"
const PathASNDB = "GeoLite2-ASN.mmdb"

type App struct {
	Cfg  config.Config
	Logg logger.Logger

	ResponseSender api.ResponseSender
	Cooker         cookie.Cooker
	Server         server.Server
	Router         router.Router

	GeoChecker *security.GeoChecker
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
	defer a.GeoChecker.Close()
	return a.Server.Run(ctx, a.Router.Run())
}

func (a *App) initModules(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		// Последовательность inits имеет значение.
		a.initConfig,
		a.initLogger,
		a.initResponseSender,
		a.initRouter,
		a.initCooker,
		a.initGeoChecker,
		a.initHandlers,
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

func (a *App) initResponseSender(ctx context.Context) error {
	a.ResponseSender = response.NewResponse()
	return nil
}

func (a *App) initGeoChecker(ctx context.Context) error {
	checker, err := security.NewGeoChecker(PathCountryDB, PathASNDB)
	if err != nil {
		return err
	}
	a.GeoChecker = checker
	return nil
}

func (a *App) initCooker(ctx context.Context) error {
	// Add default config for cookie with  Access Token
	configAccessToken := cookie.Config{
		Name:     config.COOKIE_NAME_ACCESS_TOKEN,
		Expires:  config.COOKIE_EXPIRES_ACCESS_TOKEN,
		MaxAge:   config.COOKIE_MAX_AGE_ACCESS_TOKEN,
		Path:     config.COOKIE_PATH_ACCESS_TOKEN,
		HttpOnly: config.COOKIE_HTTP_ONLY_ACCESS_TOKEN,
		Secure:   config.COOKIE_SECURE_ACCESS_TOKEN,
	}

	// Add default config for cookie with  Refresh Token
	configRefreshToken := cookie.Config{
		Name:     config.COOKIE_NAME_REFRESH_TOKEN,
		Expires:  config.COOKIE_EXPIRES_REFRESH_TOKEN,
		MaxAge:   config.COOKIE_MAX_AGE_REFRESH_TOKEN,
		Path:     config.COOKIE_PATH_REFRESH_TOKEN,
		HttpOnly: config.COOKIE_HTTP_ONLY_REFRESH_TOKEN,
		Secure:   config.COOKIE_SECURE_REFRESH_TOKEN,
	}

	a.Cooker = cookie.NewCookies(configAccessToken, configRefreshToken)
	return nil
}

func (a *App) initHandlers(ctx context.Context) error {
	// Infra services
	validater := infra.NewValidService(ctx)
	notifier := infra.NewEmailServ(ctx)
	jwtManager := jwtservice.NewJWTService()

	// Repo
	userRepo := userrepo.NewUserRepo()
	sessionRepo := userrepo.NewSessionRepo()

	// Services
	authService := service.NewAuthServ(ctx, validater, notifier, jwtManager, a.GeoChecker, userRepo, sessionRepo)
	userService := service.NewUserServ(ctx, validater, notifier, userRepo, jwtManager)

	// Handlers
	authHandler := handlers.NewAuthHandler("/auth", authService, a.ResponseSender, a.Cooker)
	userHandler := handlers.NewUserHandler("/api/v1/user", userService, a.ResponseSender)

	// Router
	a.Router.RegisterRoutes(userHandler)
	a.Router.RegisterRoutes(authHandler)
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	middlew := middleware.NewMiddlew(a.ResponseSender)
	a.Router = router.NewRouterChi(ctx, middlew)
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
