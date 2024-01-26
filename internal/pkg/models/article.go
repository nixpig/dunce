package models

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type Article struct {
	Db Dbconn
}

type ArticleData struct {
	Id        int
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,max=255"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TypeId    string    `validate:"required"`
	UserId    int       `validate:"required"`
	TagIds    string    `validate:"required"` // stored as comma separated list in db
}

type NewArticleData struct {
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,max=255"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TypeId    string    `validate:"required"`
	UserId    int       `validate:"required"`
	TagIds    string    `validate:"required"` // stored as comma separated list in db
}

func (a *Article) GetAll() (*[]ArticleData, error) {
	query := `select id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_ from article_`

	rows, err := a.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []ArticleData

	for rows.Next() {
		var article ArticleData

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a *Article) Create(newArticle NewArticleData) (*ArticleData, error) {
	validate := validator.New()

	if err := validate.Struct(newArticle); err != nil {
		return nil, err
	}

	query := `insert into article_ (title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_`

	row := a.Db.QueryRow(context.Background(), query, &newArticle.Title, &newArticle.Subtitle, &newArticle.Slug, &newArticle.Body, &newArticle.CreatedAt, &newArticle.UpdatedAt, &newArticle.TypeId, &newArticle.UserId, &newArticle.TagIds)

	var createdArticle ArticleData

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt, &createdArticle.TypeId, &createdArticle.UserId, &createdArticle.TagIds); err != nil {
		return nil, err
	}

	return &createdArticle, nil
}
