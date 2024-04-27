package articles

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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
	tagIds []int,
) Article {
	return Article{
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    tagIds,
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
	tagIds []int,
) Article {
	return Article{
		Id:        id,
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    tagIds,
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
	articleInsertQuery := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`

	row := a.db.QueryRow(context.Background(), articleInsertQuery, article.Title, article.Subtitle, article.Slug, article.Body, article.CreatedAt, article.UpdatedAt)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt); err != nil {
		return nil, err
	}

	batch := &pgx.Batch{}
	tagInsertQuery := `insert into article_tags_ (article_id_, tag_id_) values ($1, $2) returning (tag_id_)`

	for _, t := range article.TagIds {
		batch.Queue(tagInsertQuery, createdArticle.Id, t)
	}

	br := a.db.SendBatch(context.Background(), batch)

	// TODO: unwrap this once pgxmock supports batching
	if br != nil {
		defer br.Close()

		for range article.TagIds {
			var addedTagId int

			row := br.QueryRow()

			if err := row.Scan(&addedTagId); err != nil {
				return nil, err
			}

			createdArticle.TagIds = append(createdArticle.TagIds, addedTagId)
		}
	}

	return &createdArticle, nil
}

func (a ArticleData) getAll() (*[]Article, error) {
	// query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distince t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`

	rows, err := a.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var articles []Article

	for rows.Next() {
		var article Article
		var articleTagIdsConcat string

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &articleTagIdsConcat); err != nil {
			return nil, err
		}

		articleTagIds := strings.Split(articleTagIdsConcat, ",")

		for _, i := range articleTagIds {
			id, err := strconv.Atoi(i)
			if err != nil {
				return nil, err
			}

			article.TagIds = append(article.TagIds, id)
		}

		articles = append(articles, article)
	}

	return &articles, nil
}
