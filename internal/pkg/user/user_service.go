package user

import (
	"github.com/go-playground/validator/v10"
	dunce "github.com/nixpig/dunce/internal/pkg"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	data dunce.Data[UserRequest, UserResponse]
}

func NewUserService(data dunce.Data[UserRequest, UserResponse]) UserService {
	return UserService{data}
}

func (u *UserService) Create(newUser UserRequest) (*UserResponse, error) {
	validate := validator.New()

	if err := validate.Struct(newUser); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return nil, err
	}

	// TODO: check for duplicates??
	createdUser, err := u.data.Create(UserRequest{
		Username: newUser.Username,
		Email:    newUser.Email,
		Link:     newUser.Link,
		Role:     newUser.Role,
		Password: string(hashedPassword),
	})

	return createdUser, err
}

func (u *UserService) GetAll() (*[]UserResponse, error) {
	return u.data.GetAll()
}
