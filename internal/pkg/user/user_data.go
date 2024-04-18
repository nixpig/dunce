package user

import (
	"context"
	"log"

	"github.com/nixpig/dunce/internal/pkg/models"
)

type UserData struct {
	db models.Dbconn
}

func NewUserData(db models.Dbconn) UserData {
	return UserData{db}
}

func (u *UserData) Save(newUser UserNew) (*User, error) {
	query := `insert into users_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	row := u.db.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, newUser.Password)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		log.Printf("error creating user: { %s, %s, %s, %s }, %v", newUser.Username, newUser.Email, newUser.Link, newUser.Role, err)
		return nil, err
	}

	log.Printf("created user: %v", user)

	return &user, nil
}

func (u *UserData) GetAll() (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from users_`

	rows, err := u.db.Query(context.Background(), query)
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
