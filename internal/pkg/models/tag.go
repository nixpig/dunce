package models

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
