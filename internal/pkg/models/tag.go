package models

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required,max=100"`
	Slug string `json:"slug" validate:"required,max=100"`
}
