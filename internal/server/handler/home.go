package handler

import (
	"html/template"
	"net/http"
)

type HomeHandler struct {
	tmpl *template.Template
}

func NewHomeHandler(tmpl *template.Template) HomeHandler {
	return HomeHandler{tmpl: tmpl}
}

func (d HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	d.tmpl.ExecuteTemplate(w, "index.html", nil)
}
