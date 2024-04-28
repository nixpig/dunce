package articles

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/pkg/logging"
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
	Create(article *Article) (*Article, error)
	GetAll() (*[]Article, error)
	GetBySlug(slug string) (*Article, error)
	Update(article *Article) (*Article, error)
}

type ArticleData struct {
	db  db.Dbconn
	log logging.Logger
}

func NewArticleData(db db.Dbconn, log logging.Logger) ArticleData {
	return ArticleData{
		db:  db,
		log: log,
	}
}

func (a ArticleData) Create(article *Article) (*Article, error) {
	articleInsertQuery := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`
	tagInsertQuery := `insert into article_tags_ (article_id_, tag_id_) values ($1, $2) returning (tag_id_)`

	tx, err := a.db.Begin(context.Background())
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	row := tx.QueryRow(context.Background(), articleInsertQuery, article.Title, article.Subtitle, article.Slug, article.Body, article.CreatedAt, article.UpdatedAt)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt); err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	batch := &pgx.Batch{}

	for _, t := range article.TagIds {
		batch.Queue(tagInsertQuery, createdArticle.Id, t)
	}

	br := tx.SendBatch(context.Background(), batch)

	defer func() {
		var err error

		// FIXME: br isn't currently mockable, so doing this nil check to skip in tests
		if br != nil {
			err = br.Close()
		}
		if err != nil {
			a.log.Error(err.Error())
			tx.Rollback(context.Background())
		} else {
			if err := tx.Commit(context.Background()); err != nil {
				a.log.Error(err.Error())
			}
		}
	}()

	for range article.TagIds {
		// FIXME: br isn't currently mockable, so doing this nil check to skip in tests
		if br != nil {
			var addedTagId int

			row := br.QueryRow()

			if err := row.Scan(&addedTagId); err != nil {
				a.log.Error(err.Error())
				return nil, err
			}

			createdArticle.TagIds = append(createdArticle.TagIds, addedTagId)
		}
	}

	return &createdArticle, nil
}

func (a ArticleData) GetAll() (*[]Article, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`

	rows, err := a.db.Query(context.Background(), query)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	var articles []Article

	for rows.Next() {
		var article Article
		var articleTagIdsConcat string

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &articleTagIdsConcat); err != nil {
			a.log.Error(err.Error())
			return nil, err
		}

		articleTagIds := strings.Split(articleTagIdsConcat, ",")

		for _, i := range articleTagIds {
			id, err := strconv.Atoi(i)
			if err != nil {
				a.log.Error(err.Error())
				return nil, err
			}

			article.TagIds = append(article.TagIds, id)
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a ArticleData) GetBySlug(slug string) (*Article, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`

	row := a.db.QueryRow(context.Background(), query, slug)

	var article Article
	var articleTagIdsConcat string

	row.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &articleTagIdsConcat)

	articleTagIds := strings.Split(articleTagIdsConcat, ",")

	for _, i := range articleTagIds {
		tagId, err := strconv.Atoi(i)
		if err != nil {
			a.log.Error(err.Error())
			return nil, err
		}

		article.TagIds = append(article.TagIds, tagId)
	}

	return &article, nil
}

func (a ArticleData) Update(article *Article) (*Article, error) {
	updateArticleQuery := `update articles_ set title_ = $2, subtitle_ = $3, slug_ = $4, body_ = $5, created_at_ = $6, updated_at_ = $7 where id_ = $1 returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`
	deleteTagsQuery := `delete from article_tags_ where article_id_ = $1`
	updateTagsQuery := `insert into article_tags_ (article_id_, tag_id_) values ($1, $2) returning tag_id_`

	tx, err := a.db.Begin(context.Background())
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	row := tx.QueryRow(context.Background(), updateArticleQuery, &article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt)

	updatedArticle := Article{}

	if err := row.Scan(&updatedArticle.Id, &updatedArticle.Title, &updatedArticle.Subtitle, &updatedArticle.Slug, &updatedArticle.Body, &updatedArticle.CreatedAt, &updatedArticle.UpdatedAt); err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	_, err = tx.Exec(context.Background(), deleteTagsQuery, updatedArticle.Id)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	batch := &pgx.Batch{}

	for _, t := range article.TagIds {
		batch.Queue(updateTagsQuery, updatedArticle.Id, t)
	}

	br := tx.SendBatch(context.Background(), batch)

	defer func() {
		err := br.Close()
		if err != nil {
			a.log.Error(err.Error())
			tx.Rollback(context.Background())
		} else {
			if err := tx.Commit(context.Background()); err != nil {
				a.log.Error(err.Error())
			}
		}
	}()

	for range updatedArticle.TagIds {
		var updatedTagId int

		row := br.QueryRow()

		if err := row.Scan(&updatedTagId); err != nil {
			a.log.Error(err.Error())
			return nil, err
		}

		updatedArticle.TagIds = append(updatedArticle.TagIds, updatedTagId)
	}

	return &updatedArticle, nil
}
