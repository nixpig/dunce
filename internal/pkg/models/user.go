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
	Role     RoleName `json:"role_" validate:"required"`
}

type NewUser struct {
	Username string   `json:"username" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,max=100"`
	Link     string   `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	Role     RoleName `json:"role_" validate:"required"`
}

func (newUser *NewUser) CreateWithPassword(password string) (*User, error) {
	query := `insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user User

	row := database.DB.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, password)
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
