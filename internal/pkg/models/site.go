package models

type Site struct {
	Id          int      `json:"id"`
	Name        string   `json:"name" validate:"required,max=100"`
	Description string   `json:"description" validate:"required,max=255"`
	Url         string   `json:"url" validate:"required,max=255"`
	Owner       UserData `json:"owner" validate:"required"`
}
