package user

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/nixpig/dunce/pkg"
)

type UserController struct {
	log            pkg.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func NewUserController(config pkg.ControllerConfig) UserController {
	return UserController{
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
	}
}

func (u *UserController) UserLoginGet(w http.ResponseWriter, r *http.Request) {
	if err := u.templateCache["login.tmpl"].ExecuteTemplate(w, "base", nil); err != nil {
		u.log.Error(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (u *UserController) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	user := UserLogin{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	fmt.Fprintf(w, "username: '%s' | password: '%s'", user.Username, user.Password)
}

func (u *UserController) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "logout POST")
}
