package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Db Dbconn
}

type UserData struct {
	Username string   `validate:"required,min=5,max=16"`
	Email    string   `validate:"required,email,max=100"`
	Link     string   `validate:"omitempty,url,max=255"`
	Role     RoleName `validate:"required"`
}

type User struct {
	Id int `validate:"required"`
	UserData
}

type UserError []string

func NewUserError(e error) UserError {
	var userError UserError
	var ve validator.ValidationErrors

	if errors.As(e, &ve) {
		userError = make(UserError, len(ve))
		for i, fe := range ve {
			userError[i] = userError.messageForFieldError(fe)
		}
	} else {
		userError = []string{e.Error()}
	}

	return userError
}

func (u UserError) Error() string {
	return strings.Join(u, "\n")
}

func (u *UserError) messageForFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%v field is required but is empty", fe.Field())
	case "min":
		return fmt.Sprintf("%v field requires a min length of %v; length of value provided is %v", fe.Field(), fe.Param(), len(fe.Value().(string)))
	case "max":
		return fmt.Sprintf("%v field has a max length of %v; length of value provided is %v", fe.Field(), fe.Param(), len(fe.Value().(string)))
	case "email":
		return fmt.Sprintf("%v field requires an email but received %v", fe.Field(), fe.Value())
	case "url":
		return fmt.Sprintf("%v field requires a URL but received %v", fe.Field(), fe.Value())
	default:
		return "some other error"
	}
}

func (u *UserModel) Create(newUser *UserData, password string) (*User, error) {
	// TODO: interface out and inject the validator
	validate := validator.New()

	if err := validate.Struct(newUser); err != nil {
		userError := NewUserError(err)
		return nil, userError
	}

	if len(password) < 8 {
		return nil, NewUserError(errors.New("Password must be longer than 7 characters"))
	}

	// TODO: interface out and inject encryption library (will be able to test properly then!)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("unable to encrypt password: ", err)
		return nil, fmt.Errorf("unable to encrypt password")
	}

	dupeCheckQuery := `select count(id_) from users_ where username_ = $1 or email_ = $2`
	dupeRows := u.Db.QueryRow(context.Background(), dupeCheckQuery, newUser.Username, newUser.Email)

	var dupeCount int
	if err := dupeRows.Scan(&dupeCount); err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if dupeCount > 0 {
		return nil, UserError{"User already exists"}
	}

	query := `insert into users_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	var user User

	row := u.Db.QueryRow(context.Background(), query, newUser.Username, newUser.Email, newUser.Link, newUser.Role, string(hashedPassword))

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		log.Println("unable to scan row: ", err)
		return nil, fmt.Errorf("unable to scan row")
	}

	return &user, nil
}

func (u *UserModel) UpdateById(id int, user *UserData) (*User, error) {
	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		userError := NewUserError(err)
		return nil, userError
	}

	query := `update users_ set username_ = $2, email_ = $3, link_ = $4, role_ = $5 where id_ = $1 returning id_, username_, email_, link_, role_`

	row := u.Db.QueryRow(context.Background(), query, id, &user.Username, &user.Email, &user.Link, &user.Role)

	var updatedUser User

	err := row.Scan(&updatedUser.Id, &updatedUser.Username, &updatedUser.Email, &updatedUser.Link, &updatedUser.Role)

	if err == pgx.ErrNoRows {
		return nil, UserError{"User does not exist"}
	} else if err != nil {
		log.Println("unable to scan row: ", err)
		return nil, fmt.Errorf("unable to scan row")
	}

	return &updatedUser, nil
}

func (u *UserModel) GetAll() (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from users_`

	rows, err := u.Db.Query(context.Background(), query)
	if err != nil {
		log.Println("unable to execute query: ", err)
		return nil, fmt.Errorf("unable to execut query")
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			log.Println("unable to scan row: ", err)
			return nil, fmt.Errorf("unable to scan row")
		}

		users = append(users, user)
	}

	return &users, nil
}

func (u *UserModel) GetByRole(role RoleName) (*[]User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where role_ = $1`

	rows, err := u.Db.Query(context.Background(), query, role)
	if err != nil {
		log.Println("unable to execute query: ", err)
		return nil, fmt.Errorf("unable to execute query")
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
			log.Println("unable to scan row: ", err)
			return nil, fmt.Errorf("unable to scan row")
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
		log.Println("unable to scan row: ", err)
		return nil, fmt.Errorf("unable to scan row")
	}

	return &user, nil
}

func (u *UserModel) GetByUsername(username string) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where username_ = $1`

	row := u.Db.QueryRow(context.Background(), query, username)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		log.Println("unable to scan row: ", err)
		return nil, fmt.Errorf("unable to scan row")
	}

	return &user, nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `select id_, username_, email_, link_, role_ from users_ where email = $1`

	row := u.Db.QueryRow(context.Background(), query, email)

	var user User

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Link, &user.Role); err != nil {
		log.Println("unable to scan row: ", err)
		return nil, fmt.Errorf("unable to scan row")
	}

	return &user, nil
}

func (u *UserModel) Exists(s string) (bool, error) {
	query := `select id_ from users_ where username_ = $1 or email_ = $1`

	fmt.Println("check if exists: ", s)

	rows, err := u.Db.Query(context.Background(), query, s)
	if err != nil {
		log.Println("unable to execute query: ", err)
		return false, fmt.Errorf("unable to execute query")
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (u *UserModel) Delete(id int) error {
	query := `delete from users_ where id_ = $1`

	res, err := u.Db.Exec(context.Background(), query, id)
	if err != nil || res.RowsAffected() == 0 {
		log.Println("unable to execute query: ", err)
		return fmt.Errorf("unable to execute query")
	}

	return nil
}
