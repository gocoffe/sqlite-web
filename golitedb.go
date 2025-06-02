package golitedb

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/antlko/golitedb/internal/server"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed internal/db/migrations/*.sql
var embedMigrations embed.FS

func Start(dbInstance *sqlx.DB, port string) error {
	// Apply migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set migrations dialect: %w", err)
	}

	if err := goose.Up(dbInstance.Unsafe().DB, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	templates := template.Must(template.ParseFS(templatesFS, "templates/*"))

	err := server.Start(server.Config{Port: port}, dbInstance, templates)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
