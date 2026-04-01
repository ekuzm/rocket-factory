package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/ekuzm/rocket-factory/order/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a, err := app.New(ctx)
	if err != nil {
		slog.Error("order app create failed", "service", "order-service", "error", err)
		return
	}

	if err = a.Run(ctx); err != nil {
		slog.Error("order app run failed", "service", "order-service", "error", err)
		return
	}
}
