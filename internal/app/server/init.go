package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
	"gophkeep/internal/datasources"
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
	db          datasources.Datasource
}

func NewApp(ctx context.Context, cfg config.ConfigProvider) (*App, error) {
	db, err := postgresdb.NewPostgresStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	app := &App{cfgProvider: cfg, db: db}

	router := gin.Default()

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
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_ = app.db.Close()
	return app.server.Shutdown(shutdownCtx)

}

func (app *App) registerHandlers(router *gin.Engine, h handlers.HandlersProvider, m middleware.MiddlewareProvider) error {
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	apiPath := router.Group("/api")

	apiPath.POST(`/login`, h.LoginUser)
	apiPath.POST(`/auth`, h.RegisterUser)
	apiPath.GET(`/ping`, m.VerifyJWT, h.Ping)

	secretsPath := apiPath.Group("/secrets", m.VerifyJWT)

	secretsPath.POST(`/`, h.AddSecret)
	secretsPath.GET(`/`, h.GetAllSecrets)
	secretsPath.GET(`/:secretID`, h.GetSecret)
	secretsPath.PUT(`/:secretID`, h.UpdateSecret)
	secretsPath.DELETE(`/:secretID`, h.DeleteSecret)

	return nil
}
