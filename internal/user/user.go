package user

import "time"

type UserLogin struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserNew struct {
	Username string `validate:"required"`
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
