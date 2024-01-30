package models

import (
	"context"
	"fmt"
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
	TypeId    int       `validate:"required"`
	TypeName  string    `validate:"required"`
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
	TypeId    int       `validate:"required"`
	UserId    int       `validate:"required"`
	TagIds    string    `validate:"required"` // stored as comma separated list in db
}

func (a *Article) GetBySlug(slug string) (*ArticleData, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.name_ from articles_ a inner join types_ t on a.type_id_ = t.id_ where a.slug_ = $1`

	row := a.Db.QueryRow(context.Background(), query, slug)

	var article ArticleData

	if err := row.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.TypeName); err != nil {
		return nil, err
	}

	return &article, nil
}

func (a *Article) GetAll() (*[]ArticleData, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.name_ from articles_ a inner join types_ t on a.type_id_ = t.id_`

	rows, err := a.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []ArticleData

	for rows.Next() {
		var article ArticleData

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.TypeName); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a *Article) GetByTypeName(typeName string) (*[]ArticleData, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.name_ from articles_ a inner join types_ t on a.type_id_ = t.id_ where t.name_ = $1`

	rows, err := a.Db.Query(context.Background(), query, typeName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []ArticleData

	for rows.Next() {
		var article ArticleData

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.TypeName); err != nil {
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

	query := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_`

	fmt.Println(newArticle.Title, newArticle.Subtitle, newArticle.Slug, newArticle.Body, newArticle.CreatedAt, newArticle.UpdatedAt, newArticle.TypeId, newArticle.UserId, newArticle.TagIds)

	row := a.Db.QueryRow(context.Background(), query, &newArticle.Title, &newArticle.Subtitle, &newArticle.Slug, &newArticle.Body, &newArticle.CreatedAt, &newArticle.UpdatedAt, &newArticle.TypeId, &newArticle.UserId, &newArticle.TagIds)

	var createdArticle ArticleData

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt, &createdArticle.TypeId, &createdArticle.UserId, &createdArticle.TagIds); err != nil {
		return nil, err
	}

	return &createdArticle, nil
}
