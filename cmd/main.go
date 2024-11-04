package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/geomodular/go-htmx-todo/internal/server"
)

const host = ":8080"

func main() {
	ctx := context.Background()

	slog.Info("starting server", "host", host)
	if err := server.Run(ctx, host); err != nil {
		slog.Error("failed running server", "msg", err)
		os.Exit(1)
	}
	os.Exit(0)
}
