package models

import "time"

type Article struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" validate:"required,max=255"`
	Subtitle  string    `json:"subtitle" validate:"required,max=255"`
	Slug      string    `json:"slug" validate:"required,max=255"`
	Body      string    `json:"body" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	TypeId    string    `json:"type_id" validate:"required"`
	UserId    int       `json:"user_id" validate:"required"`
	TagIds    []int     `json:"tag_ids" validate:"required"` // stored as comma separated list in db
}
