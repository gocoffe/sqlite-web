package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/antlko/golitedb/internal/db"
	"github.com/antlko/golitedb/internal/server/handler"
	"github.com/jmoiron/sqlx"
)

func Start(cfg Config, dbInstance *sqlx.DB, templates *template.Template) error {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	userRepo := db.NewUserRepo(dbInstance)
	tableRepo := db.NewTablesRepo(dbInstance)

	loginHandler := handler.NewLoginHandler(templates, &userRepo)
	dashboardHandler := handler.NewDashboardHandler(templates, dbInstance, &tableRepo)
	homeHandler := handler.NewHomeHandler(templates)

	http.HandleFunc("/", homeHandler.GetHome)

	http.HandleFunc("POST /login", loginHandler.POSTHandle)
	http.HandleFunc("POST /logout", loginHandler.POSTLogoutHandle)
	http.HandleFunc("GET /login", loginHandler.GETHandle)
	http.HandleFunc("POST /change-password", loginHandler.POSTChangePasswordHandle)

	http.HandleFunc("GET /dashboard", dashboardHandler.GetDashboard)
	http.HandleFunc("GET /dashboard/table", dashboardHandler.GetTableRows)
	http.HandleFunc("POST /dashboard/console/exec", dashboardHandler.ExecSQL)

	slog.InfoContext(context.Background(), "DB UI Server started!")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
