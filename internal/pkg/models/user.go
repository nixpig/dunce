package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=100"`
	Link     string `json:"link" validate:"required,max=255"` // e.g. twitter.com/nixpig
	RoleId   int    `json:"role_id" validate:"required"`
}
