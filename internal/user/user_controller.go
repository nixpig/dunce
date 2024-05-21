package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
)

type UserController struct {
	service        UserService
	log            logging.Logger
	templateCache  templates.TemplateCache
	sessionManager session.SessionManager
	csrfToken      func(r *http.Request) string
	errorHandlers  errors.ErrorHandlers
}

type UserControllerConfig struct {
	Log            logging.Logger
	TemplateCache  templates.TemplateCache
	SessionManager session.SessionManager
	CsrfToken      func(*http.Request) string
	ErrorHandlers  errors.ErrorHandlers
}

type UserView struct {
	Message         string
	User            *UserResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

type UsersView struct {
	Message         string
	Users           *[]UserResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

type UserLoginView struct {
	Message         string
	CsrfToken       string
	IsAuthenticated bool
}

type UserCreateView struct {
	CsrfToken       string
	IsAuthenticated bool
}

func NewUserController(
	service UserService,
	config UserControllerConfig,
) UserController {
	return UserController{
		service:        service,
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
		csrfToken:      config.CsrfToken,
		errorHandlers:  config.ErrorHandlers,
	}
}

func (u *UserController) UserLoginGet(w http.ResponseWriter, r *http.Request) {
	if u.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
		return
	}

	if err := u.templateCache["pages/admin/login.tmpl"].ExecuteTemplate(w, "admin", UserLoginView{
		Message:         u.sessionManager.PopString(r.Context(), session.SESSION_KEY_MESSAGE),
		CsrfToken:       u.csrfToken(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
		u.errorHandlers.InternalServerError(w, r)
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
		u.sessionManager.Put(r.Context(), session.SESSION_KEY_MESSAGE, "Login failed.")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if err := u.sessionManager.RenewToken(r.Context()); err != nil {
		u.log.Error(err.Error())
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	u.sessionManager.Put(r.Context(), session.LOGGED_IN_USERNAME, username)

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (u *UserController) UserLogoutPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	u.sessionManager.Remove(r.Context(), session.LOGGED_IN_USERNAME)
	u.sessionManager.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		"You've been logged out.",
	)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (u *UserController) CreateUserGet(w http.ResponseWriter, r *http.Request) {
	if err := u.templateCache["pages/admin/new-user.tmpl"].ExecuteTemplate(w, "admin", UserCreateView{
		CsrfToken:       u.csrfToken(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (u *UserController) CreateUserPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := UserNewRequestDto{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	createdUser, err := u.service.Create(&user)
	if err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}

	u.sessionManager.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		fmt.Sprintf("Created user '%s'.", createdUser.Username),
	)

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (u *UserController) UsersGet(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.GetAll()
	if err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}

	message := u.sessionManager.PopString(r.Context(), session.SESSION_KEY_MESSAGE)

	if err := u.templateCache["pages/admin/users.tmpl"].ExecuteTemplate(w, "admin", UsersView{
		Message:         message,
		Users:           users,
		CsrfToken:       u.csrfToken(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (u *UserController) UserGet(w http.ResponseWriter, r *http.Request) {
	user, err := u.service.GetByAttribute("username", r.PathValue("slug"))
	if err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}

	if err := u.templateCache["pages/admin/user.tmpl"].ExecuteTemplate(w, "admin", UserView{
		Message:         "",
		User:            user,
		CsrfToken:       u.csrfToken(r),
		IsAuthenticated: u.IsAuthenticated(r),
	}); err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (u *UserController) DeleteUserPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	username := r.FormValue("username")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		u.errorHandlers.BadRequest(w, r)
		return
	}

	if err := u.service.DeleteById(uint(id)); err != nil {
		u.errorHandlers.InternalServerError(w, r)
		return
	}

	u.sessionManager.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		fmt.Sprintf("Deleted user '%s'.", username),
	)

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (u *UserController) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(session.IS_LOGGED_IN_CONTEXT_KEY).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
