package handler

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gocoffe/sqlite-web/internal/db"
	"github.com/gocoffe/sqlite-web/internal/jwt"
	"github.com/jmoiron/sqlx"
)

type DashboardHandler struct {
	tmpl       *template.Template
	dbInstance *sqlx.DB
	tablesRepo *db.TablesRepo
}

func NewDashboardHandler(tmpl *template.Template, dbInstance *sqlx.DB, tablesRepo *db.TablesRepo, j *jwt.Authorizer) DashboardHandler {
	return DashboardHandler{tmpl: tmpl, dbInstance: dbInstance, tablesRepo: tablesRepo}
}

func (d DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	tables, err := d.tablesRepo.GetAllDBTables()
	if err != nil {
		d.tmpl.ExecuteTemplate(w, "dashboard.html", map[string]interface{}{
			"Error": err.Error(),
		})
	}

	d.tmpl.ExecuteTemplate(w, "dashboard.html", map[string]interface{}{
		"Tables": tables,
	})
}

func (d DashboardHandler) GetTableRows(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("name")
	if tableName == "" {
		http.Error(w, "Missing table name", http.StatusBadRequest)
		return
	}

	query := "SELECT * FROM " + tableName + " LIMIT 50"
	rows, err := d.dbInstance.Queryx(query)
	if err != nil {
		http.Error(w, "Failed to fetch table data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err == nil {
			results = append(results, row)
		}
	}

	writeJSON(w, results)
}

func (d DashboardHandler) ExecSQL(w http.ResponseWriter, r *http.Request) {
	slog.Info("ExecSQL request handle")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.FormValue("query")
	if query == "" {
		http.Error(w, "Query is empty", http.StatusBadRequest)
		return
	}

	// Basic inspection
	trimmed := strings.TrimSpace(strings.ToUpper(query))
	isSelect := strings.HasPrefix(trimmed, "SELECT")

	if isSelect {
		rows, err := d.dbInstance.Queryx(query)
		if err != nil {
			http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			row := make(map[string]interface{})
			if err := rows.MapScan(row); err == nil {
				results = append(results, row)
			}
		}
		writeJSON(w, results)
	} else {
		res, err := d.dbInstance.Exec(query)
		if err != nil {
			http.Error(w, "Execution error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		affected, _ := res.RowsAffected()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Query executed",
			"rows_affected": affected,
		})
	}
}

// Хелпер для JSON ответа
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
