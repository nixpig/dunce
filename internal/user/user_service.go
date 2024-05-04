package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
)

type UserService struct {
	repo     pkg.Repository[User, UserNew]
	validate *validator.Validate
	log      pkg.Logger
}

func NewUserService(
	repo pkg.Repository[User, UserNew],
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
	return u.repo.Create(user)
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
