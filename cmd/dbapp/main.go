package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"github.com/gocoffe/sqlite-web/internal/app/dbapp"
	"github.com/gocoffe/sqlite-web/internal/db"
	"github.com/gocoffe/sqlite-web/internal/jwt"
	"github.com/gocoffe/sqlite-web/internal/server"
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

	err := dbapp.Start(readFlags(cfg))
	if err != nil {
		slog.Error(err.Error())
	}
}

func readFlags(cfg dbapp.Config) dbapp.Config {
	dbFilePath := flag.String("f", "sqlite.db", "SQLiteDB file path, current path by default")
	serverPort := flag.String("p", "5577", "Server Port")
	appName := flag.String("app", "sqlite-web", "Application Name")
	secretKey := flag.String("secret-key", "secret-key", "Secret Key")

	flag.Parse()

	return dbapp.Config{
		ApplicationName: notEmptyString(pointerToString(appName), cfg.ApplicationName),
		DB: db.Config{
			FilePath:     notEmptyString(pointerToString(dbFilePath), cfg.DB.FilePath),
			MaxIdleConns: cfg.DB.MaxIdleConns,
			MaxOpenConns: cfg.DB.MaxOpenConns,
			LogLevel:     cfg.DB.LogLevel,
		},
		Server: server.Config{
			Port: notEmptyString(pointerToString(serverPort), cfg.Server.Port),
		},
		Authorizer: jwt.Config{
			JwtSecretKey:        notEmptyString(pointerToString(secretKey), cfg.Authorizer.JwtSecretKey),
			JwtAccessTokenHours: cfg.Authorizer.JwtAccessTokenHours,
		},
	}
}

func notEmptyString(v1, v2 string) string {
	if v1 != "" {
		return v1
	}
	return v2
}

func pointerToString(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
