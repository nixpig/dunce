package user

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(user *UserNew) (*User, error)
	DeleteById(id int) error
	Exists(username string) (bool, error)
	GetAll() (*[]User, error)
	GetByAttribute(attr, value string) (*User, error)
	Update(user *User) (*User, error)
	LoginWithUsernamePassword(username, password string) error
}

type UserService struct {
	repo     IUserRepository
	validate *validator.Validate
}

func NewUserService(
	repo IUserRepository,
	validate *validator.Validate,
) UserService {
	return UserService{
		repo:     repo,
		validate: validate,
	}
}

func (u UserService) Create(user *UserNew) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}

	return u.repo.Create(&UserNew{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
	})
}

func (u UserService) GetAll() (*[]User, error) {
	return u.repo.GetAll()
}

func (u UserService) GetByAttribute(attr, value string) (*User, error) {
	return u.repo.GetByAttribute(attr, value)
}

func (u UserService) Update(user *User) (*User, error) {
	return u.repo.Update(user)
}

func (u UserService) DeleteById(id int) error {
	return u.repo.DeleteById(id)
}

func (u UserService) LoginWithUsernamePassword(username, password string) error {
	hashedPassword, err := u.repo.GetPasswordByUsername(username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (u UserService) Exists(username string) (bool, error) {
	exists, err := u.repo.Exists(username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u UserService) IsAuthenticatedMiddleware(sessionManager *scs.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := sessionManager.GetString(r.Context(), pkg.LOGGED_IN_USERNAME)

		if len(username) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := u.Exists(username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), pkg.IsLoggedInContextKey, true)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	})
}
