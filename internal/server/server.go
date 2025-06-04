package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/antlko/golitedb/internal/db"
	"github.com/antlko/golitedb/internal/jwt"
	"github.com/antlko/golitedb/internal/server/handler"
	"github.com/antlko/golitedb/internal/server/middleware"
	"github.com/jmoiron/sqlx"
)

const (
	defaultHashCost = 8
)

func Start(cfg Config, dbInstance *sqlx.DB, templates *template.Template, authorizer jwt.Authorizer) error {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	hasher := jwt.NewPasswordHasher(defaultHashCost)

	userRepo := db.NewUserRepo(dbInstance)
	tableRepo := db.NewTablesRepo(dbInstance)

	loginHandler := handler.NewLoginHandler(templates, &userRepo, &authorizer, &hasher)
	dashboardHandler := handler.NewDashboardHandler(templates, dbInstance, &tableRepo, &authorizer)
	homeHandler := handler.NewHomeHandler(templates)

	http.HandleFunc("/", homeHandler.GetHome)

	http.HandleFunc("POST /login", loginHandler.POSTLogin)
	http.HandleFunc("POST /logout", loginHandler.POSTLogoutHandle)
	http.HandleFunc("GET /login", loginHandler.GETHandle)
	http.HandleFunc("POST /change-password", loginHandler.POSTChangePasswordHandle)

	http.HandleFunc("GET /dashboard", middleware.NewAuthMiddleware(&authorizer, dashboardHandler.GetDashboard).Handle)
	http.HandleFunc("GET /dashboard/table", dashboardHandler.GetTableRows)
	http.HandleFunc("POST /dashboard/console/exec", dashboardHandler.ExecSQL)

	slog.InfoContext(context.Background(), "DB WebUI Server started!")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
