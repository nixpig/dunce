package user

import "time"

type User struct {
	Id        uint   `validate:"omitempty"`
	Username  string `validate:"required"`
	Email     string `validate:"required"`
	Password  string `validate:"required"`
	CreatedAt time.Time
}

type UserLoginRequestDto struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type UserNewRequestDto struct {
	Username string `validate:"required"`
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type UserResponseDto struct {
	Id       uint   `validate:"required"`
	Username string `validate:"required"`
	Email    string `validate:"required"`
}
