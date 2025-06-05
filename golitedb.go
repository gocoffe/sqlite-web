package golitedb

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/antlko/golitedb/internal/jwt"
	"github.com/antlko/golitedb/internal/server"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

const (
	defaultPort                  = "5566"
	defaultSessionExpTimeMinutes = 60
	defaultSessionSecretKey      = "secret-key"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed internal/db/migrations/*.sql
var embedMigrations embed.FS

//go:embed static/*
var staticFiles embed.FS

type Params struct {
	DBInstance *sqlx.DB
	AppConfig  Config
}

type Config struct {
	ServerPort        string
	SessionSecretKey  string
	SessionTokenHours int64
}

func Start(params Params) error {
	if params.AppConfig.ServerPort == "" {
		params.AppConfig.ServerPort = defaultPort
	}
	if params.AppConfig.SessionSecretKey == "" {
		params.AppConfig.SessionSecretKey = defaultSessionSecretKey
	}
	if params.AppConfig.SessionTokenHours == 0 {
		params.AppConfig.SessionTokenHours = defaultSessionExpTimeMinutes
	}

	// Apply migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set migrations dialect: %w", err)
	}

	if err := goose.Up(params.DBInstance.Unsafe().DB, "internal/db/migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	templates := template.Must(template.ParseFS(templatesFS, "templates/*"))

	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	jwtAuthorizer := jwt.NewAuthorizer(jwt.Config{
		JwtSecretKey:        "",
		JwtAccessTokenHours: 0,
	})

	err := server.Start(
		server.Config{
			Port: params.AppConfig.ServerPort,
		},
		params.DBInstance,
		templates,
		jwtAuthorizer,
	)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
