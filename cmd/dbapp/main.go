package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gocoffe/sqlite-web/internal/app/dbapp"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load .env file", "error", err)
	}

	var cfg dbapp.Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		slog.Error("Failed to process configuration", "error", err)
		return
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := dbapp.Start(cfg)
	if err != nil {
		slog.Error(err.Error())
	}
}
