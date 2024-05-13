package middleware

import "net/http"

func stripSlashMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isRoot := r.URL.Path == "/" || len(r.URL.Path) == 0

		if isRoot {
			next(w, r)
			return
		}

		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
			http.Redirect(w, r, r.URL.Path, 301)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewStripSlashMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return stripSlashMiddleware(next)
	}
}
