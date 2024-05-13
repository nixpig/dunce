package middleware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/nixpig/dunce/internal/user"
	"github.com/nixpig/dunce/pkg"
)

func NewAuthenticatedMiddleware(userService user.UserService, sessionManager *scs.SessionManager) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return AuthenticatedMiddleware(userService, sessionManager, next)
	}
}

func AuthenticatedMiddleware(userService user.UserService, sessionManager *scs.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := sessionManager.GetString(r.Context(), pkg.LOGGED_IN_USERNAME)

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
			ctx := context.WithValue(r.Context(), pkg.IS_LOGGED_IN_CONTEXT_KEY, true)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	})
}
