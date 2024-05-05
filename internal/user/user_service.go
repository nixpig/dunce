package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo     UserRepository
	validate *validator.Validate
	log      pkg.Logger
}

func NewUserService(
	repo UserRepository,
	validate *validator.Validate,
	log pkg.Logger,
) UserService {
	return UserService{
		repo:     repo,
		validate: validate,
		log:      log,
	}
}

func (u UserService) Create(user *UserNew) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		u.log.Error(err.Error())
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
		u.log.Error(err.Error())
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}
