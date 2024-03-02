package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mrz1836/go-sanitize"
)

type ArticleModel struct {
	Db Dbconn
}

type ArticleData struct {
	Title     string    `validate:"required,max=255"`
	Subtitle  string    `validate:"required,max=255"`
	Slug      string    `validate:"required,slug,max=255"`
	Body      string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	TypeId    int       `validate:"required"`
	UserId    int       `validate:"required"`
	TagIds    string    `validate:"required"` // stored as comma separated list in db
}

type Article struct {
	Id int
	ArticleData

	Type Type `validate:"required"`
	// Tags []TagData `validate:"required"`
	// User UserData  `validate:"required"`
}

func (a *ArticleModel) GetById(id int) (*Article, error) {
	fmt.Println("get by id")
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.id_, t.name_, t.slug_ from articles_ a inner join types_ t on a.type_id_ = t.id_ where a.id_ = $1`

	row := a.Db.QueryRow(context.Background(), query, id)

	var article Article

	if err := row.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.Type.Id, &article.Type.Name, &article.Type.Slug); err != nil {
		fmt.Println("errored: ", err)
		return nil, err
	}

	return &article, nil
}

func (a *ArticleModel) GetBySlug(slug string) (*Article, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.id_, t.name_, t.slug_ from articles_ a inner join types_ t on a.type_id_ = t.id_ where a.slug_ = $1`

	row := a.Db.QueryRow(context.Background(), query, slug)

	var article Article

	if err := row.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.Type.Id, &article.Type.Name, &article.Type.Slug); err != nil {
		return nil, err
	}

	return &article, nil
}

func (a *ArticleModel) GetAll() (*[]Article, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.id_, t.name_, t.slug_ from articles_ a inner join types_ t on a.type_id_ = t.id_`

	rows, err := a.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.Type.Id, &article.Type.Name, &article.Type.Slug); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a *ArticleModel) GetByTypeName(typeName string) (*[]Article, error) {
	query := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, a.type_id_, a.user_id_, a.tag_ids_, t.id_, t.name_, t.slug_, from articles_ a inner join types_ t on a.type_id_ = t.id_ where t.name_ = $1`

	rows, err := a.Db.Query(context.Background(), query, typeName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.Id, &article.Title, &article.Subtitle, &article.Slug, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.TypeId, &article.UserId, &article.TagIds, &article.Type.Id, &article.Type.Name, &article.Type.Slug); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (a *ArticleModel) Create(newArticle ArticleData) (*Article, error) {
	sanitisedData := ArticleData{
		Title:     sanitize.XSS(newArticle.Title),
		Subtitle:  sanitize.XSS(newArticle.Subtitle),
		Slug:      sanitize.PathName(newArticle.Slug),
		Body:      sanitize.XSS(newArticle.Body),
		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
		TypeId:    newArticle.TypeId,
		UserId:    newArticle.UserId,
		TagIds:    sanitize.XSS(newArticle.TagIds),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedData); err != nil {
		return nil, err
	}

	checkDuplicatesQuery := `select count(*) from articles_ where slug_ = $1`

	var duplicateCount int
	dupRow := a.Db.QueryRow(context.Background(), checkDuplicatesQuery, &sanitisedData.Slug)
	if err := dupRow.Scan(&duplicateCount); err != nil {
		return nil, err
	}

	if duplicateCount > 0 {
		return nil, fmt.Errorf("duplicate slug")
	}

	query := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_`

	fmt.Println(sanitisedData.Title, sanitisedData.Subtitle, sanitisedData.Slug, sanitisedData.Body, sanitisedData.CreatedAt, sanitisedData.UpdatedAt, sanitisedData.TypeId, sanitisedData.UserId, sanitisedData.TagIds)

	row := a.Db.QueryRow(context.Background(), query, &sanitisedData.Title, &sanitisedData.Subtitle, &sanitisedData.Slug, &sanitisedData.Body, &sanitisedData.CreatedAt, &sanitisedData.UpdatedAt, &sanitisedData.TypeId, &sanitisedData.UserId, &sanitisedData.TagIds)

	var createdArticle Article

	if err := row.Scan(&createdArticle.Id, &createdArticle.Title, &createdArticle.Subtitle, &createdArticle.Slug, &createdArticle.Body, &createdArticle.CreatedAt, &createdArticle.UpdatedAt, &createdArticle.TypeId, &createdArticle.UserId, &createdArticle.TagIds); err != nil {
		return nil, err
	}

	return &createdArticle, nil
}

func (a *ArticleModel) UpdateById(id int, updateArticle Article) (*Article, error) {
	sanitisedData := Article{
		ArticleData: ArticleData{
			Title:     sanitize.XSS(updateArticle.Title),
			Subtitle:  sanitize.XSS(updateArticle.Subtitle),
			Slug:      sanitize.PathName(updateArticle.Slug),
			Body:      sanitize.XSS(updateArticle.Body),
			CreatedAt: updateArticle.CreatedAt,
			UpdatedAt: updateArticle.UpdatedAt,
			TypeId:    updateArticle.TypeId,
			UserId:    updateArticle.UserId,
			TagIds:    sanitize.XSS(updateArticle.TagIds),
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedData); err != nil {
		return nil, err
	}

	checkDuplicatesQuery := `select count(*) from articles_ where slug_ = $2 and id_ <> $1`

	var duplicateCount int
	dupRow := a.Db.QueryRow(context.Background(), checkDuplicatesQuery, id, &sanitisedData.Slug)
	if err := dupRow.Scan(&duplicateCount); err != nil {
		return nil, err
	}

	if duplicateCount > 0 {
		return nil, fmt.Errorf("duplicate slug")
	}

	query := `update articles_ set title_ = $2, subtitle_ = $3, slug_ = $4, body_ = $5, created_at_ = $6, updated_at_ = $7, type_id_ = $8, user_id_ = $9, tag_ids_ = $10 where id_ = $1 returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_, type_id_, user_id_, tag_ids_`

	fmt.Println(sanitisedData.Title, sanitisedData.Subtitle, sanitisedData.Slug, sanitisedData.Body, sanitisedData.CreatedAt, sanitisedData.UpdatedAt, sanitisedData.TypeId, sanitisedData.UserId, sanitisedData.TagIds)

	row := a.Db.QueryRow(context.Background(), query, id, &sanitisedData.Title, &sanitisedData.Subtitle, &sanitisedData.Slug, &sanitisedData.Body, &sanitisedData.CreatedAt, &sanitisedData.UpdatedAt, &sanitisedData.TypeId, &sanitisedData.UserId, &sanitisedData.TagIds)

	var updatedArticleData Article

	if err := row.Scan(&updatedArticleData.Id, &updatedArticleData.Title, &updatedArticleData.Subtitle, &updatedArticleData.Slug, &updatedArticleData.Body, &updatedArticleData.CreatedAt, &updatedArticleData.UpdatedAt, &updatedArticleData.TypeId, &updatedArticleData.UserId, &updatedArticleData.TagIds); err != nil {
		return nil, err
	}

	return &updatedArticleData, nil
}

func (a *ArticleModel) DeleteById(id int) error {
	query := `delete from articles_ where id_ = $1`

	res, err := a.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("did not delete any rows")
	}

	return nil
}
