package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gophkeep/internal/app/server"
	"gophkeep/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	cfg.Logger().Info("Db conn string", zap.String("db-conn-str", cfg.GetConfig().DBConnStr))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer func() {
		cfg.Logger().Info("Received an interrupt, shutting down...")
		stop()
	}()

	app, err := server.NewApp(ctx, cfg)
	if err != nil {
		panic(err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("app run failed: %w", err)
		}
		cfg.Logger().Info("app stopped")
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		return app.Stop(ctx)
	})

	if err := g.Wait(); err != nil {
		cfg.Logger().Fatal(err.Error())
	}

}
