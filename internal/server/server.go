package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gocoffe/sqlite-web/internal/db"
	"github.com/gocoffe/sqlite-web/internal/jwt"
	"github.com/gocoffe/sqlite-web/internal/server/handler"
	"github.com/gocoffe/sqlite-web/internal/server/middleware"
	"github.com/jmoiron/sqlx"
)

const (
	defaultHashCost = 8
)

func Start(cfg Config, dbInstance *sqlx.DB, templates *template.Template, authorizer jwt.Authorizer) error {
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
	http.HandleFunc("GET /dashboard/table", middleware.NewAuthMiddleware(&authorizer, dashboardHandler.GetTableRows).Handle)
	http.HandleFunc("POST /dashboard/console/exec", middleware.NewAuthMiddleware(&authorizer, dashboardHandler.ExecSQL).Handle)

	slog.InfoContext(context.Background(), "DB WebUI Server started!")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
