package models

type Login struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=255"`
}
