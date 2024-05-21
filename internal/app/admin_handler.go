package app

import (
	"net/http"

	"github.com/nixpig/dunce/pkg/session"
)

// func registerAdminRoutes(mux *http.ServeMux) *http.ServeMux {
//
// }

func adminRootHandler(appConfig AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appConfig.SessionManager.Exists(
			r.Context(),
			string(session.IS_LOGGED_IN_CONTEXT_KEY),
		) {
			http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
	}
}
