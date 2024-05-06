package pkg

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

func Protected(sessionManager *scs.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sessionManager.Exists(r.Context(), LOGGED_IN_USERNAME) {
			sessionManager.Put(r.Context(), SESSION_KEY_MESSAGE, "You are not logged in.")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfHttpHandler := nosurf.New(next)
		csrfHttpHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
		})

		csrfHttpHandler.ServeHTTP(w, r)
	})
}
