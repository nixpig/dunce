package article

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/tag"
)

type ArticleRepository interface {
	DeleteById(id int) error
	Create(article *ArticleNew) (*Article, error)
	GetAll() (*[]Article, error)
	GetManyByAttribute(attr, value string) (*[]Article, error)
	GetByAttribute(attr, value string) (*Article, error)
	Update(article *UpdateArticle) (*Article, error)
}

type articlePostgresRepository struct {
	db db.Dbconn
}

func NewArticlePostgresRepository(db db.Dbconn) articlePostgresRepository {
	return articlePostgresRepository{
		db: db,
	}
}

func (a articlePostgresRepository) DeleteById(id int) error {
	query := `delete from articles_ a where a.id_ = $1`

	_, err := a.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (a articlePostgresRepository) Create(article *ArticleNew) (*Article, error) {
	articleInsertQuery := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`
	tagInsertQuery := `with tags as (select id_, name_, slug_ from tags_ where id_ = $2), article_tags as (insert into article_tags_ (article_id_, tag_id_) values ($1, $2)) select id_, name_, slug_ from tags`

	tx, err := a.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(context.Background(), articleInsertQuery, article.Title, article.Subtitle, article.Slug, article.Body, article.CreatedAt, article.UpdatedAt)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt); err != nil {
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
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	for range article.TagIds {
		// FIXME: br isn't currently mockable, so doing this nil check to skip in tests
		if br != nil {
			var addedTag tag.Tag

			row := br.QueryRow()

			if err := row.Scan(&addedTag.Id, &addedTag.Name, &addedTag.Slug); err != nil {
				return nil, err
			}

			createdArticle.Tags = append(createdArticle.Tags, addedTag)
		}
	}

	return &createdArticle, nil
}

func (a articlePostgresRepository) GetAll() (*[]Article, error) {
	articlesQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`
	tagsQuery := `select id_, name_, slug_ from tags_`

	tagRows, err := a.db.Query(context.Background(), tagsQuery)
	if err != nil {
		return nil, err
	}

	defer tagRows.Close()

	tagList := map[int]tag.Tag{}

	for tagRows.Next() {
		var singleTag tag.Tag

		if err := tagRows.Scan(&singleTag.Id, &singleTag.Name, &singleTag.Slug); err != nil {
			return nil, err
		}

		tagList[singleTag.Id] = singleTag
	}

	rows, err := a.db.Query(context.Background(), articlesQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article
		var articleTagIdsConcat string

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &articleTagIdsConcat); err != nil {
			return nil, err
		}

		articleTagIds := strings.Split(articleTagIdsConcat, ",")

		articleTags := make([]tag.Tag, len(articleTagIds))

		for index, articleTagId := range articleTagIds {
			id, err := strconv.Atoi(articleTagId)
			if err != nil {
				return nil, err
			}

			articleTags[index] = tagList[id]

		}

		article.Tags = articleTags

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a articlePostgresRepository) GetManyByAttribute(attr, value string) (*[]Article, error) {
	var articleQuery string
	var err error

	switch attr {
	case "tagSlug":
		articleQuery = `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_ from articles_ a inner join article_tags_ at on a.id_ = at.article_id_ inner join tags_ t on at.tag_id_ = t.id_ where t.slug_ = $1`

	default:
		return nil, errors.New("unsupported attribute")
	}

	rows, err := a.db.Query(context.Background(), articleQuery, value)
	if err != nil {
		return nil, err
	}

	var articles []Article

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	// TODO: fetch and attach tags to article results

	return &articles, nil
}

func (a articlePostgresRepository) GetByAttribute(attr, value string) (*Article, error) {
	var articleQuery string

	switch attr {
	case "slug":
		articleQuery = `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`
	default:
		return nil, errors.New("invalid attribute")
	}

	row := a.db.QueryRow(context.Background(), articleQuery, value)

	var article Article
	var articleTagIdsConcat string

	if err := row.Scan(
		&article.Id,
		&article.Title,
		&article.Subtitle,
		&article.Slug,
		&article.Body,
		&article.CreatedAt,
		&article.UpdatedAt,
		&articleTagIdsConcat,
	); err != nil {
		return nil, err
	}

	// TODO: hopefully pgx supports scanning postgres arrays so don't need to perform multiple queries or this funky string concat
	tagsQuery := strings.Join(
		[]string{
			`select id_, name_, slug_ from tags_ where id_ = `,
			strings.ReplaceAll(articleTagIdsConcat, ",", " or id_ = "),
		}, "")

	rows, err := a.db.Query(context.Background(), tagsQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag tag.Tag
		if err := rows.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}

		article.Tags = append(article.Tags, tag)
	}

	return &article, nil
}

func (a articlePostgresRepository) Update(article *UpdateArticle) (*Article, error) {
	updateArticleQuery := `update articles_ set title_ = $2, subtitle_ = $3, slug_ = $4, body_ = $5, created_at_ = $6, updated_at_ = $7 where id_ = $1 returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`
	deleteTagsQuery := `delete from article_tags_ where article_id_ = $1`
	updateTagsQuery := `insert into article_tags_ (article_id_, tag_id_) values ($1, $2) returning tag_id_`
	tagsQuery := `select id_, name_, slug_ from tags_`

	tagsRows, err := a.db.Query(context.Background(), tagsQuery)
	if err != nil {
		return nil, err
	}

	tagList := map[int]tag.Tag{}

	for tagsRows.Next() {
		var tag tag.Tag

		if err := tagsRows.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}

		tagList[tag.Id] = tag
	}

	tx, err := a.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(context.Background(), updateArticleQuery, &article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt)

	updatedArticle := Article{}

	if err := row.Scan(&updatedArticle.Id, &updatedArticle.Title, &updatedArticle.Subtitle, &updatedArticle.Slug, &updatedArticle.Body, &updatedArticle.CreatedAt, &updatedArticle.UpdatedAt); err != nil {
		return nil, err
	}

	_, err = tx.Exec(context.Background(), deleteTagsQuery, updatedArticle.Id)
	if err != nil {
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
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	// FIXME: this doesn't seem right... check out 11 lines below
	for range updatedArticle.Tags {
		var updatedTagId int

		row := br.QueryRow()

		if err := row.Scan(&updatedTagId); err != nil {
			return nil, err
		}

		updatedArticle.Tags = append(updatedArticle.Tags, tagList[updatedTagId])
	}

	return &updatedArticle, nil
}
