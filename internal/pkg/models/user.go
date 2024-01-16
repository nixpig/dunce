package models

import (
	"context"

	"github.com/nixpig/bloggor/internal/pkg/database"
)

type User struct {
	Id       int      `json:"id"`
	Username string   `json:"username" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,max=100"`
	Link     string   `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	Role     RoleName `json:"role" validate:"required"`
}

type NewUser struct {
	Username string   `json:"username" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,max=100"`
	Link     string   `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	Password string   `json:"password" validate:"required,max=255"`
	Role     RoleName `json:"role" validate:"required"`
}

func CreateUser(newUser *NewUser) (*User, error) {
	query := `insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user User

	row := database.DB.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, newUser.Password)

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsers() (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from user_`

	rows, err := database.DB.Query(context.Background(), query)
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

func GetUsersByRole(role RoleName) (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where role_ = $1`

	rows, err := database.DB.Query(context.Background(), query, role)
	if err != nil {
		return nil, err
	}

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

func GetUserById(id int) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from user_ where id_ = $1`

	var user User

	row := database.DB.QueryRow(context.Background(), query, id)

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
