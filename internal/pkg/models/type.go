package models

type Type struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required,max=255"`
	Template string `json:"template" validate:"required,max=255"`
	Slug     string `json:"slug" validate:"required,max=255"`
}
