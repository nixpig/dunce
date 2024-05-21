package app

import "net/http"

func publicRootHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.Error(w, "Not Found", 404)
			return
		}

		next(w, r)
	}
}
