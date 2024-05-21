package middleware

import (
	"net/http"

	"github.com/nixpig/dunce/pkg/session"
)

func NewProtectedMiddleware(sessionManager session.SessionManager) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return ProtectedMiddleware(sessionManager, next)
	}
}

func ProtectedMiddleware(sessionManager session.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sessionManager.Exists(r.Context(), session.LOGGED_IN_USERNAME) {
			sessionManager.Put(r.Context(), session.SESSION_KEY_MESSAGE, "You are not logged in.")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
