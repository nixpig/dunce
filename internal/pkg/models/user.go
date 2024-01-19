package models

import (
	"context"
)

type User struct {
	Db Dbconn
}

type UserData struct {
	Id       int      `json:"id"`
	Username string   `json:"username" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,max=100"`
	Link     string   `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	Role     RoleName `json:"role" validate:"required"`
}

type NewUserData struct {
	Username string   `json:"username" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,max=100"`
	Link     string   `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	Password string   `json:"password" validate:"required,max=255"`
	Role     RoleName `json:"role" validate:"required"`
}

func (u *User) Create(newUser *NewUserData) (*UserData, error) {
	query := `insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user UserData

	row := u.Db.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, newUser.Password)

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
