package user

import (
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/pkg"
)

type UserRepository struct {
	db  db.Dbconn
	log pkg.Logger
}

func NewUserRepository(db db.Dbconn, log pkg.Logger) UserRepository {
	return UserRepository{
		db:  db,
		log: log,
	}
}

func (u UserRepository) Authenticate(user *User) (int, error) {
	return 0, nil
}

func (u UserRepository) Create(user *UserNew) (*User, error) {
	return nil, nil
}

func (u UserRepository) DeleteById(id int) error {
	return nil
}

func (u UserRepository) Exists(user *User) (bool, error) {
	return false, nil
}

func (u UserRepository) GetAll() (*[]User, error) {
	return &[]User{}, nil
}

func (u UserRepository) GetByAttribute(attr, value string) (*User, error) {
	return nil, nil
}

func (u UserRepository) Update(user *User) (*User, error) {
	return nil, nil
}
