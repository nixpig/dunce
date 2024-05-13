package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/nixpig/dunce/pkg"
)

func NewProtectedMiddleware(sessionManager *scs.SessionManager) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return ProtectedMiddleware(sessionManager, next)
	}
}

func ProtectedMiddleware(sessionManager *scs.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sessionManager.Exists(r.Context(), pkg.LOGGED_IN_USERNAME) {
			sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, "You are not logged in.")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
