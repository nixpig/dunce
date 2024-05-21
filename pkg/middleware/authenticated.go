package middleware

import (
	"context"
	"net/http"

	"github.com/nixpig/dunce/internal/user"
	"github.com/nixpig/dunce/pkg/session"
)

func NewAuthenticatedMiddleware(userService user.UserService, sessionManager session.SessionManager, sessionKey string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return AuthenticatedMiddleware(userService, sessionManager, sessionKey, next)
	}
}

func AuthenticatedMiddleware(userService user.UserService, sessionManager session.SessionManager, sessionKey string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := sessionManager.GetString(r.Context(), sessionKey)

		if len(username) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := userService.Exists(username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, true)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	})
}
