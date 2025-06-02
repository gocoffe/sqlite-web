package handler

import (
	"html/template"
	"log/slog"
	"net/http"
)

func respondError(tmpl *template.Template, w http.ResponseWriter, file string, error string) {
	w.WriteHeader(http.StatusUnauthorized)
	if err := tmpl.ExecuteTemplate(w, file, map[string]string{
		"error": error,
	}); err != nil {
		slog.Error("Templ Execution", "error", err)
	}
}
