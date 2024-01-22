package models

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Db Dbconn
}

type UserData struct {
	Id       int      ``
	Username string   `validate:"required,max=100"`
	Email    string   `validate:"required,max=100"`
	Link     string   `validate:"required,max=255"`
	Role     RoleName `validate:"required"`
}

type NewUserData struct {
	Username string   `json:"username" validate:"required,min=5,max=100"`
	Email    string   `json:"email" validate:"required,email,max=100"`
	Link     string   `json:"link" validate:"required,url,max=255"`
	Password string   `json:"password" validate:"required,min=8,max=255"`
	Role     RoleName `json:"role" validate:"required"`
}

func (u *User) Create(newUser *NewUserData) (*UserData, error) {
	validate := validator.New()

	if err := validate.Struct(newUser); err != nil {
		return nil, err
	}

	// usernameTaken, err := u.Exists(newUser.Username)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if usernameTaken {
	// 	return nil, fmt.Errorf("username already taken")
	// }
	//
	// emailTaken, err := u.Exists(newUser.Email)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if emailTaken {
	// 	return nil, fmt.Errorf("email already taken")
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return nil, err
	}

	query := `insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user UserData

	row := u.Db.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, hashedPassword)

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetAll() (*[]UserData, error) {
	query := `select id_, username_, email_, link_, role_ from user_`

	rows, err := u.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []UserData

	for rows.Next() {
		var user UserData

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u *User) GetByRole(role RoleName) (*[]UserData, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where role_ = $1`

	rows, err := u.Db.Query(context.Background(), query, role)
	if err != nil {
		return nil, err
	}

	var users []UserData

	for rows.Next() {
		var user UserData

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u *User) GetById(id int) (*UserData, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where id_ = $1`

	var user UserData

	row := u.Db.QueryRow(context.Background(), query, id)

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByUsername(username string) (*UserData, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where username_ = $1`

	row := u.Db.QueryRow(context.Background(), query, username)

	var user UserData

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetByEmail(email string) (*UserData, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where email = $1`

	row := u.Db.QueryRow(context.Background(), query, email)

	var user UserData

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Exists(s string) (bool, error) {
	query := `select id_ from user_ where username_ = $1 or email_ = $1`

	fmt.Println("check if exists: ", s)

	rows, err := u.Db.Query(context.Background(), query, s)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}
