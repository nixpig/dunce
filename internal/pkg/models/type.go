package models

type Type struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Template string `json:"template"`
	Slug     string `json:"slug"`
}
