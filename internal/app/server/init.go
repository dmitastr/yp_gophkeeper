package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
	"gophkeep/internal/datasources/postgresdb"
	"gophkeep/internal/domain/service"
	"gophkeep/internal/presentation/handlers"
	"gophkeep/internal/presentation/middleware"
)

type IApp interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	cfgProvider config.ConfigProvider
	server      *http.Server
}

func NewApp(ctx context.Context, cfg config.ConfigProvider) (*App, error) {
	app := &App{cfgProvider: cfg}
	router := gin.Default()

	// db := datasources.NewDummyDS(cfg)
	db, err := postgresdb.NewPostgresStorage(context.TODO(), cfg)
	if err != nil {
		return nil, err
	}

	serviceProvider := service.NewService(cfg, db)
	handlersProvider := handlers.NewHandlersProvider(cfg, serviceProvider)

	middlewareProvider := middleware.NewMiddlewareProvider(cfg, serviceProvider)

	err = app.registerHandlers(router, handlersProvider, middlewareProvider)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:              cfg.GetConfig().Address,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           router,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	app.server = srv

	return app, nil
}

func (app *App) Run() error {
	app.cfgProvider.Logger().Info("Starting gophkeep app")
	return app.server.ListenAndServe()
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) registerHandlers(router *gin.Engine, h handlers.HandlersProvider, m middleware.MiddlewareProvider) error {
	apiPath := router.Group("/api")
	apiPath.POST(`/login`, h.LoginUser)
	apiPath.POST(`/auth`, h.RegisterUser)
	apiPath.GET(`/ping`, m.VerifyJWT, h.Ping)

	secretsPath := apiPath.Group("/secrets", m.VerifyJWT)

	secretsPath.POST(`/`, h.AddSecret)
	secretsPath.GET(`/`, h.GetAllSecrets)
	secretsPath.GET(`/:secretID`, h.GetSecret)
	secretsPath.POST(`/passwords`, h.AddPassword)

	return nil
}
