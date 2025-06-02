package handler

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/antlko/golitedb/internal/db"
)

type Login struct {
	tmpl     *template.Template
	userRepo *db.UserRepo
}

func NewLoginHandler(tmpl *template.Template, userRepo *db.UserRepo) Login {
	return Login{
		tmpl:     tmpl,
		userRepo: userRepo,
	}
}

func (l Login) POSTHandle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	login := r.FormValue("login")
	password := r.FormValue("password")

	slog.Info("Login request", "login", login)

	dbUser, err := l.userRepo.GetByLogin(ctx, login)
	if err != nil {
		println(err.Error())
		respondError(l.tmpl, w, "login.html", "Incorrect Credentials!")
		return
	}

	// TODO: Dummy password check, improve with hashing password
	if password != dbUser.Password {
		respondError(l.tmpl, w, "login.html", "Incorrect Credentials!")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "some my sess token", // TODO: create normal JWT token
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

func (l Login) GETHandle(w http.ResponseWriter, r *http.Request) {
	if err := l.tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		slog.Error("Templ Execution", "error", err)
	}
}

func (l Login) POSTLogoutHandle(w http.ResponseWriter, r *http.Request) {
	println("logout handle")
	// Clear the session token by setting an expired cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (l Login) POSTChangePasswordHandle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldPass := r.FormValue("old_password")
	newPass := r.FormValue("new_password")

	// TODO: in future should be get user from the token
	login := "admin"

	user, err := l.userRepo.GetByLogin(ctx, login)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if user.Password != oldPass {
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	err = l.userRepo.UpdatePassword(ctx, login, newPass)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
