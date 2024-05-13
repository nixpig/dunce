package user

import (
	"context"
	"errors"

	"github.com/nixpig/dunce/db"
)

type UserRepository interface {
	Create(user *User) (*User, error)
	DeleteById(id uint) error
	Exists(username string) (bool, error)
	GetAll() (*[]User, error)
	GetByAttribute(attr, value string) (*User, error)
	GetPasswordByUsername(username string) (string, error)
	Update(user *User) (*User, error)
}

type userPostgresRepository struct {
	db db.Dbconn
}

func NewUserPostgresRepository(db db.Dbconn) userPostgresRepository {
	return userPostgresRepository{
		db: db,
	}
}

func (u userPostgresRepository) Create(user *User) (*User, error) {
	query := `insert into users_ (username_, email_, password_) values ($1, $2, $3) returning id_, username_, email_`

	var createdUser User

	row := u.db.QueryRow(context.Background(), query, user.Username, user.Email, user.Password)

	if err := row.Scan(&createdUser.Id, &createdUser.Username, &createdUser.Email); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (u userPostgresRepository) DeleteById(id uint) error {
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

func (u userPostgresRepository) Exists(username string) (bool, error) {
	query := `select exists(select true from users_ where username_ = $1)`

	row := u.db.QueryRow(context.Background(), query, username)

	var exists bool

	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (u userPostgresRepository) GetAll() (*[]User, error) {
	query := `select id_, username_, email_ from users_`

	rows, err := u.db.Query(context.Background(), query)
	if err != nil {
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

func (u userPostgresRepository) GetByAttribute(attr, value string) (*User, error) {
	var query string

	switch attr {
	case "username":
		query = `select id_, username_, email_ from users_ where username_ = $1`
	default:
		err := errors.New("attribute not supported")
		return nil, err
	}

	row := u.db.QueryRow(context.Background(), query, value)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userPostgresRepository) GetPasswordByUsername(username string) (string, error) {
	query := `select password_ from users_ where username_ = $1`

	row := u.db.QueryRow(context.Background(), query, username)

	var password string

	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (u userPostgresRepository) Update(user *User) (*User, error) {
	return nil, nil
}
