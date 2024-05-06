package user

import (
	"context"
	"errors"
	"fmt"

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
	query := `insert into users_ (username_, email_, password_) values ($1, $2, $3) returning id_, username_, email_`

	row := u.db.QueryRow(context.Background(), query, user.Username, user.Email, user.Password)

	createdUser := User{}

	if err := row.Scan(&createdUser.Id, &createdUser.Username, &createdUser.Email); err != nil {
		fmt.Println("OOPS!!", err)
		return nil, err
	}

	return &createdUser, nil
}

func (u UserRepository) DeleteById(id int) error {
	query := `delete from users_ where id_ = $1`

	res, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("no user deleted")
	}

	return nil
}

func (u UserRepository) Exists(username string) (bool, error) {
	query := `select exists(select true from users_ where username_ = $1)`

	row := u.db.QueryRow(context.Background(), query, username)

	var exists bool

	if err := row.Scan(&exists); err != nil {
		u.log.Error(err.Error())
		return false, err
	}

	return exists, nil
}

func (u UserRepository) GetAll() (*[]User, error) {
	query := `select id_, username_, email_ from users_`

	rows, err := u.db.Query(context.Background(), query)
	if err != nil {
		u.log.Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u UserRepository) GetByAttribute(attr, value string) (*User, error) {
	var query string

	switch attr {
	case "username":
		query = `select id_, username_, email_ from users_ where username_ = $1`
	default:
		err := errors.New("attribute not supported")
		u.log.Error(err.Error())
		return nil, err
	}

	row := u.db.QueryRow(context.Background(), query, value)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) GetPasswordByUsername(username string) (string, error) {
	query := `select password_ from users_ where username_ = $1`

	row := u.db.QueryRow(context.Background(), query, username)

	var password string

	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (u UserRepository) Update(user *User) (*User, error) {
	return nil, nil
}
