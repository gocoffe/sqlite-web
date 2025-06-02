package golitedb

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/antlko/golitedb/internal/server"
	"github.com/jmoiron/sqlx"
)

//go:embed templates/*
var templatesFS embed.FS

func Start(dbInstance *sqlx.DB, port string) error {
	templates := template.Must(template.ParseFS(templatesFS, "templates/*"))

	err := server.Start(server.Config{Port: port}, dbInstance, templates)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
