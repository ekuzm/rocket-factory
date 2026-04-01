package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/ekuzm/rocket-factory/inventory/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a, err := app.New(ctx)
	if err != nil {
		slog.Error("Failed to create app", "service", "inventory-service", "error", err)
		return
	}

	if err := a.Run(ctx); err != nil {
		slog.Error("Failed to run app", "service", "inventory-service", "error", err)
	}
}
