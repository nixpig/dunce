package articles

import (
	"context"
	"time"

	"github.com/nixpig/dunce/db"
)

type Article struct {
	Id        int
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,min=2,max=50"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	Tags      []int     `validate:"required"`
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
	tags []int,
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
	tags []int,
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

type ArticleDataInterface interface {
	create(article *Article) (*Article, error)
	getAll() (*[]Article, error)
}

type ArticleData struct {
	db db.Dbconn
}

func NewArticleData(db db.Dbconn) ArticleData {
	return ArticleData{db}
}

func (a ArticleData) create(article *Article) (*Article, error) {
	query := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`

	row := a.db.QueryRow(context.Background(), query, article.Title, article.Subtitle, article.Slug, article.Body, article.CreatedAt, article.UpdatedAt)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt); err != nil {
		return nil, err
	}

	return &createdArticle, nil
}

func (a ArticleData) getAll() (*[]Article, error) {
	// query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, t.id_, t.name_, t.slug_ from articles_ a inner join types_ t on a.type_id_ = t.id_`

	return nil, nil
}
