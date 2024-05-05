package user

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/nixpig/dunce/pkg"
)

type UserController struct {
	service        pkg.Service[User, UserNew]
	log            pkg.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func NewUserController(service pkg.Service[User, UserNew], config pkg.ControllerConfig) UserController {
	return UserController{
		service:        service,
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

func (u *UserController) CreateUserGet(w http.ResponseWriter, r *http.Request) {
	if err := u.templateCache["new-user.tmpl"].ExecuteTemplate(w, "base", nil); err != nil {
		u.log.Error(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func (u *UserController) CreateUserPost(w http.ResponseWriter, r *http.Request) {
	user := UserNew{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	createdUser, err := u.service.Create(&user)
	if err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	u.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Created user '%s'.", createdUser.Username))

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (u *UserController) UsersGet(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.GetAll()
	if err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := u.sessionManager.PopString(r.Context(), pkg.SESSION_KEY_MESSAGE)

	if err := u.templateCache["users.tmpl"].ExecuteTemplate(w, "base", struct {
		Message string
		Users   *[]User
	}{
		Message: message,
		Users:   users,
	}); err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (u *UserController) UserGet(w http.ResponseWriter, r *http.Request) {
	user, err := u.service.GetByAttribute("username", r.PathValue("slug"))
	if err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := u.templateCache["user.tmpl"].ExecuteTemplate(w, "base", struct {
		Message string
		User    *User
	}{
		Message: "",
		User:    user,
	}); err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (u *UserController) DeleteUserPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := u.service.DeleteById(id); err != nil {
		u.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	u.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Deleted user '%s'.", username))

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
