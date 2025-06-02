package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/antlko/golitedb/internal/app/dbapp"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func Start() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load .env file", "error", err)
	}

	// Load application configuration
	var cfg dbapp.Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		slog.Error("Failed to process configuration", "error", err)
		return
	}

	println(fmt.Sprintf("%+v", cfg))

	dbapp.Start(cfg)
}
