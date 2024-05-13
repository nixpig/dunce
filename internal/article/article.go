package article

import (
	"time"

	"github.com/nixpig/dunce/internal/tag"
)

type Article struct {
	Id        int       `validate:"omitempty"`
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	Tags      []tag.Tag `validate:"required"`
}

type ArticleNewRequestDto struct {
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TagIds    []int     `validate:"required"`
}

type ArticleUpdateRequestDto struct {
	Id        int       `validate:"omitempty"`
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TagIds    []int     `validate:"required"`
}

type UpdateArticle struct {
	Id        int       `validate:"omitempty"`
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TagIds    []int     `validate:"required"`
}

type ArticleNew struct {
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TagIds    []int     `validate:"required"`
}

type ArticleResponseDto struct {
	Id        int
	Title     string
	Subtitle  string
	Slug      string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []tag.Tag
}

type ArticleTag struct {
	Id        int `validate:"required"`
	ArticleId int `validate:"required"`
	TagId     int `validate:"required"`
}
