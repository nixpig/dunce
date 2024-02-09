package models

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/mrz1836/go-sanitize"
)

type Tag struct {
	Db Dbconn
}

type TagData struct {
	Id   int
	Name string
	Slug string
}

type NewTagData struct {
	Name string `validate:"required,max=100"`
	Slug string `validate:"required,slug,max=100"`
}

type UpdateTagData struct {
	Id   int    `validate:"required,number"`
	Name string `validate:"required,alphanumunicode,max=100"`
	Slug string `validate:"required,slug,max=100"`
}

func (t *Tag) Create(newTag NewTagData) (*TagData, error) {
	sanitisedTagData := NewTagData{
		Name: sanitize.XSS(newTag.Name),
		Slug: sanitize.PathName(newTag.Slug),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedTagData); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	// TODO: check if tag name or slug already exists
	if t.IsDuplicate(sanitisedTagData) {
		return nil, fmt.Errorf("duplicate tag: %v", sanitisedTagData)
	}

	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var tag TagData

	row := t.Db.QueryRow(context.Background(), query, &sanitisedTagData.Name, &sanitisedTagData.Slug)

	if err := row.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (t *Tag) Delete(id int) error {
	query := `delete from tags_ where id_ = $1`

	res, err := t.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}

func (t *Tag) GetAll() (*[]TagData, error) {
	query := `select id_, name_, slug_ from tags_`

	rows, err := t.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []TagData

	for rows.Next() {
		var tag TagData

		if err := rows.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return &tags, nil
}

func (t *Tag) IsDuplicate(tag NewTagData) bool {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	var rowCount int
	row := t.Db.QueryRow(context.Background(), query, &tag.Name, &tag.Slug)
	row.Scan(&rowCount)

	return rowCount > 0
}

func ValidateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
