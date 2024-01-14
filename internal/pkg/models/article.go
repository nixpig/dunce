package models

import "time"

type Article struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Slug      string    `json:"slug"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TypeId    string    `json:"type_id"`
	UserId    int       `json:"user_id"`
	TagIds    []int     `json:"tag_ids"`
}
