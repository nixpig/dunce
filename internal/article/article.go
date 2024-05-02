package article

import (
	"time"

	"github.com/nixpig/dunce/internal/tag"
)

type Article struct {
	Id        int
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	Tags      []tag.Tag `validate:"required"`
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

type ArticleTag struct {
	Id        int
	ArticleId int
	TagId     int
}

func NewArticle(
	title string,
	subtitle string,
	slug string,
	body string,
	createdAt time.Time,
	updatedAt time.Time,
	tags []tag.Tag,
) Article {
	return Article{
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags:      tags,
	}
}

func NewArticleWithId(
	id int,
	title string,
	subtitle string,
	slug string,
	body string,
	createdAt time.Time,
	updatedAt time.Time,
	tags []tag.Tag,
) Article {
	return Article{
		Id:        id,
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags:      tags,
	}
}
