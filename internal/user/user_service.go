package user

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user *UserNewRequestDto) (*UserResponseDto, error)
	DeleteById(id uint) error
	Exists(username string) (bool, error)
	GetAll() (*[]UserResponseDto, error)
	GetByAttribute(attr, value string) (*UserResponseDto, error)
	Update(user *User) (*UserResponseDto, error)
	LoginWithUsernamePassword(username, password string) error
}

type UserServiceImpl struct {
	repo     UserRepository
	validate *validator.Validate
}

func NewUserService(
	repo UserRepository,
	validate *validator.Validate,
) UserServiceImpl {
	return UserServiceImpl{
		repo:     repo,
		validate: validate,
	}
}

func (u UserServiceImpl) Create(user *UserNewRequestDto) (*UserResponseDto, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}

	userToCreate := User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	if err := u.validate.Struct(userToCreate); err != nil {
		return nil, err
	}

	createdUser, err := u.repo.Create(&userToCreate)
	if err != nil {
		return nil, err
	}

	return &UserResponseDto{
		Id:       createdUser.Id,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}, nil
}

func (u UserServiceImpl) GetAll() (*[]UserResponseDto, error) {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	allUsers := make([]UserResponseDto, len(*users))

	for index, user := range *users {
		allUsers[index] = UserResponseDto{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		}
	}

	return &allUsers, nil
}

func (u UserServiceImpl) GetByAttribute(attr, value string) (*UserResponseDto, error) {
	user, err := u.repo.GetByAttribute(attr, value)
	if err != nil {
		return nil, err
	}

	return &UserResponseDto{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u UserServiceImpl) Update(user *User) (*UserResponseDto, error) {
	if err := u.validate.Struct(user); err != nil {
		return nil, err
	}

	user, err := u.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return &UserResponseDto{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u UserServiceImpl) DeleteById(id uint) error {
	return u.repo.DeleteById(id)
}

func (u UserServiceImpl) LoginWithUsernamePassword(username, password string) error {
	hashedPassword, err := u.repo.GetPasswordByUsername(username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (u UserServiceImpl) Exists(username string) (bool, error) {
	exists, err := u.repo.Exists(username)
	if err != nil {
		return false, err
	}

	return exists, nil
}
