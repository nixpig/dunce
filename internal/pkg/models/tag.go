package models

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mrz1836/go-sanitize"
)

type TagModel struct {
	Db Dbconn
}

type TagData struct {
	Name string `validate:"required,max=100"`
	Slug string `validate:"required,slug,max=100"`
}

type Tag struct {
	Id int
	TagData
}

func (t *TagModel) Create(newTag TagData) (*Tag, error) {
	sanitisedTagData := TagData{
		Name: sanitize.XSS(newTag.Name),
		Slug: sanitize.PathName(newTag.Slug),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedTagData); err != nil {
		log.Errorf("failed validating: %v", err)
		return nil, err.(validator.ValidationErrors)
	}

	checkDuplicatesQuery := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	var duplicateCount int

	duplicateRow := t.Db.QueryRow(context.Background(), checkDuplicatesQuery, &sanitisedTagData.Name, &sanitisedTagData.Slug)
	if err := duplicateRow.Scan(&duplicateCount); err != nil {
		log.Errorf("failed scanning duplicate: %v", err)
		return nil, err
	}

	if duplicateCount > 0 {
		log.Error("duplicate tag")
		return nil, fmt.Errorf("duplicate tag: '%s' '%s'", sanitisedTagData.Name, sanitisedTagData.Slug)
	}

	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var tag Tag

	row := t.Db.QueryRow(context.Background(), query, &sanitisedTagData.Name, &sanitisedTagData.Slug)

	if err := row.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
		log.Errorf("failed scanning tag: %v", err)
		return nil, err
	}

	return &tag, nil
}

func (t *TagModel) Delete(id int) error {
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

func (t *TagModel) UpdateById(id int, tag TagData) (*Tag, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	sanitisedTagData := TagData{
		Name: sanitize.XSS(tag.Name),
		Slug: sanitize.PathName(tag.Slug),
	}

	if err := validate.Struct(sanitisedTagData); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	var duplicateCount int

	duplicateQuery := `select count(*) from tags_ where (name_ = $2 or slug_ = $3) and id_ <> $1`

	duplicateRows := t.Db.QueryRow(context.Background(), duplicateQuery, id, &sanitisedTagData.Name, &sanitisedTagData.Slug)
	duplicateRows.Scan(&duplicateCount)

	if duplicateCount > 0 {
		return nil, fmt.Errorf("a tag with these attributes already exists")
	}

	query := `update tags_ set name_ = $2, slug_ = $3 where id_ = $1 returning id_, name_, slug_`

	row := t.Db.QueryRow(context.Background(), query, id, &sanitisedTagData.Name, &sanitisedTagData.Slug)

	var updatedTag Tag

	if err := row.Scan(&updatedTag.Id, &updatedTag.Name, &updatedTag.Slug); err != nil {
		return nil, err
	}

	return &updatedTag, nil
}

func (t *TagModel) GetById(id int) (*Tag, error) {
	query := `select id_, name_, slug_ from tags_ where id_ = $1`

	row := t.Db.QueryRow(context.Background(), query, id)

	var tag Tag

	if err := row.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (t *TagModel) GetAll() (*[]Tag, error) {
	query := `select id_, name_, slug_ from tags_`

	rows, err := t.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag

		if err := rows.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return &tags, nil
}

func ValidateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9\\-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
