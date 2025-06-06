package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gocoffe/sqlite-web/internal/jwt"
)

type Auth struct {
	next       func(w http.ResponseWriter, r *http.Request)
	authorizer *jwt.Authorizer
}

func NewAuthMiddleware(authorizer *jwt.Authorizer, next func(w http.ResponseWriter, r *http.Request)) Auth {
	return Auth{
		authorizer: authorizer,
		next:       next,
	}
}

func (a Auth) Handle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ok, identity, err := a.authorizer.Validate(cookie.Value)
	if err != nil || !ok {
		http.Error(w, "User not authorized", http.StatusUnauthorized)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "X-User", identity))

	slog.Info("user is valid", "user", identity)

	a.next(w, r)
}
