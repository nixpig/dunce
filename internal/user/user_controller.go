package user

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/nixpig/dunce/pkg"
)

type UserController struct {
	service        IUserService
	log            pkg.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func NewUserController(service IUserService, config pkg.ControllerConfig) UserController {
	return UserController{
		service:        service,
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
	}
}

func (u *UserController) UserLoginGet(w http.ResponseWriter, r *http.Request) {
	if u.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
		return
	}

	if err := u.templateCache["pages/admin/admin-login.tmpl"].ExecuteTemplate(w, "admin", struct {
		Message         string
		CsrfToken       string
		IsAuthenticated bool
	}{
		Message:         u.sessionManager.PopString(r.Context(), pkg.SESSION_KEY_MESSAGE),
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
		u.log.Error(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (u *UserController) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if err := u.service.LoginWithUsernamePassword(
		username,
		password,
	); err != nil {
		u.log.Error(err.Error())
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	if err := u.sessionManager.RenewToken(r.Context()); err != nil {
		u.log.Error(err.Error())
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	u.sessionManager.Put(r.Context(), pkg.LOGGED_IN_USERNAME, username)

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (u *UserController) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	u.sessionManager.Remove(r.Context(), pkg.LOGGED_IN_USERNAME)
	u.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, "You've been logged out.")

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (u *UserController) CreateUserGet(w http.ResponseWriter, r *http.Request) {
	if err := u.templateCache["pages/admin/admin-new-user.tmpl"].ExecuteTemplate(w, "admin", struct {
		CsrfToken       string
		IsAuthenticated bool
	}{
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
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

	if err := u.templateCache["pages/admin/admin-users.tmpl"].ExecuteTemplate(w, "admin", struct {
		Message         string
		Users           *[]User
		CsrfToken       string
		IsAuthenticated bool
	}{
		Message:         message,
		Users:           users,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: u.IsAuthenticated(r),
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

	if err := u.templateCache["pages/admin/admin-user.tmpl"].ExecuteTemplate(w, "admin", struct {
		Message         string
		User            *User
		CsrfToken       string
		IsAuthenticated bool
	}{
		Message:         "",
		User:            user,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: u.IsAuthenticated(r),
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

func (u *UserController) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(pkg.IsLoggedInContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
