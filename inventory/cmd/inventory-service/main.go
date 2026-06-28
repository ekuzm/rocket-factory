package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	_ "go.uber.org/automaxprocs"

	"github.com/ekuzm/rocket-factory/inventory/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a, err := app.New(ctx)
	if err != nil {
		slog.Error("inventory app create failed", "error", err)
		return
	}

	if err := a.Run(ctx); err != nil {
		slog.Error("inventory app run failed", "error", err)
	}
}
