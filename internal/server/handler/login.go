package handler

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/gocoffe/sqlite-web/internal/db"
	"github.com/gocoffe/sqlite-web/internal/jwt"
)

type Login struct {
	tmpl       *template.Template
	userRepo   *db.UserRepo
	authorizer *jwt.Authorizer
	hasher     *jwt.PasswordHasher
}

func NewLoginHandler(tmpl *template.Template, userRepo *db.UserRepo, authorizer *jwt.Authorizer, hasher *jwt.PasswordHasher) Login {
	return Login{
		tmpl:       tmpl,
		userRepo:   userRepo,
		authorizer: authorizer,
		hasher:     hasher,
	}
}

func (l Login) POSTLogin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	login := r.FormValue("login")
	password := r.FormValue("password")

	slog.Info("Login request", "login", login)

	dbUser, err := l.userRepo.GetByLogin(ctx, login)
	if err != nil {
		respondError(l.tmpl, w, "login.html", "Incorrect Credentials!")
		return
	}

	if dbUser.Login == "admin" && dbUser.Password == "admin" {
		// Skipping. Can be improved in the future, but now needed only for the first activation
	} else {
		if !l.hasher.VerifyPassword(password, dbUser.Password) {
			respondError(l.tmpl, w, "login.html", "Incorrect Credentials!")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	sessionToken, err := l.authorizer.CreateToken(login)
	if err != nil {
		respondError(l.tmpl, w, "login.html", "Internal Server Error!")
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
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

	sessionTokenCookie, err := r.Cookie("session_token")
	if err != nil || sessionTokenCookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldPass := r.FormValue("old_password")
	newPass := r.FormValue("new_password")

	ok, login, err := l.authorizer.Validate(sessionTokenCookie.Value)
	if err != nil || !ok {
		http.Error(w, "User invalid", http.StatusBadRequest)
		return
	}

	user, err := l.userRepo.GetByLogin(ctx, login)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if user.Password != oldPass {
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := l.hasher.HashPassword(newPass)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = l.userRepo.UpdatePassword(ctx, login, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
