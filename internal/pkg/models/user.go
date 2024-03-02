package models

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Db Dbconn
}

type UserData struct {
	Username string   `validate:"required,min=5,max=100"`
	Email    string   `validate:"required,email,max=100"`
	Link     string   `validate:"omitempty,url,max=255"`
	Role     RoleName `validate:"required"`
}

type User struct {
	Id int `validate:"required"`
	UserData
}

func (u *UserModel) Create(newUser *UserData, password string) (*User, error) {
	validate := validator.New()

	if err := validate.Struct(newUser); err != nil {
		return nil, fmt.Errorf("invalid user: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt password: %v", err)
	}

	query := `insert into users_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user User

	row := u.Db.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, string(hashedPassword))

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, fmt.Errorf("unable to scan row: %v", err)
	}

	return &user, nil
}

func (u *UserModel) Update(id int, user *UserData) (*User, error) {
	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return nil, err
	}

	query := `update users_ set username_ = $2, email_ = $3, link_ = $4, role_ = $5 where id_ = $1 returning id_, username_, email_, link_, role_`

	row := u.Db.QueryRow(context.Background(), query, id, &user.Username, &user.Email, &user.Link, &user.Role)

	var updatedUser User

	if err := row.Scan(&updatedUser.Id, &updatedUser.Username, &updatedUser.Email, &updatedUser.Link, &updatedUser.Role); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (u *UserModel) GetAll() (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from users_`

	rows, err := u.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u *UserModel) GetByRole(role RoleName) (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where role_ = $1`

	rows, err := u.Db.Query(context.Background(), query, role)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u *UserModel) GetById(id int) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where id_ = $1`

	var user User

	row := u.Db.QueryRow(context.Background(), query, id)

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) GetByUsername(username string) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where username_ = $1`

	row := u.Db.QueryRow(context.Background(), query, username)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where email = $1`

	row := u.Db.QueryRow(context.Background(), query, email)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) Exists(s string) (bool, error) {
	query := `select id_ from users_ where username_ = $1 or email_ = $1`

	fmt.Println("check if exists: ", s)

	rows, err := u.Db.Query(context.Background(), query, s)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (u *UserModel) Delete(id int) error {
	query := `delete from users_ where id_ = $1`

	res, err := u.Db.Exec(context.Background(), query, id)
	if err != nil || res.RowsAffected() == 0 {
		return err
	}

	return nil
}
