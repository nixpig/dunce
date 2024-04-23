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
	UserId    int       `validate:"required"`
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
	userId int,
) Article {
	return Article{
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserId:    userId,
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
	userId int,
) Article {
	return Article{
		Id:        id,
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserId:    userId,
	}
}

type ArticleDataInterface interface {
	create(article *Article) (*Article, error)
}

type ArticleData struct {
	db db.Dbconn
}

func NewArticleData(db db.Dbconn) ArticleData {
	return ArticleData{db}
}

func (a ArticleData) create(article *Article) (*Article, error) {
	query := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_, user_id_) values ($1, $2, $3, $4, $5, $6, $7) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, user_id_`

	row := a.db.QueryRow(context.Background(), query, article.Title, article.Subtitle, article.Slug, article.Body, article.CreatedAt, article.UpdatedAt, article.UserId)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt, &createdArticle.UserId); err != nil {
		return nil, err
	}

	return &createdArticle, nil
}
