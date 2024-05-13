package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func NewNoSurfMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return NoSurfMiddleware(next)
	}
}

func NoSurfMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
