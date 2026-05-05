package main

import (
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/commands"
	"gophkeep/internal/logger"
)

func main() {
	log := logger.NewLogger()
	root := commands.NewCmd(agent.RootDeps{
		Logger: log,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer func() {
		log.Info("Received an interrupt, shutting down agent...")
		stop()
	}()

	g, _ := errgroup.WithContext(ctx)

	// Agent app goroutine
	g.Go(func() error {
		if err := root.ExecuteContext(ctx); err != nil {
			return fmt.Errorf("gophkeeper agent error: %w", err)
		}
		log.Info("gophkeeper stopped")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Info("exit reason: %v", zap.Error(err))
	}
}
