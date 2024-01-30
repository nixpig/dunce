package models

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Tag struct {
	Db Dbconn
}

type TagData struct {
	Id   int
	Name string `validate:"required,max=100"`
	Slug string `validate:"required,max=100"`
}

type NewTagData struct {
	Name string `validate:"required,max=100"`
	Slug string `validate:"required,max=100"`
}

type UpdateTagData struct {
	Id   int    `validate:"required"`
	Name string `validate:"required,max=100"`
	Slug string `validate:"required,max=100"`
}

func (t *Tag) Create(newTag NewTagData) (*TagData, error) {

	validate := validator.New()

	if err := validate.Struct(newTag); err != nil {
		return nil, err
	}

	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var tag TagData

	row := t.Db.QueryRow(context.Background(), query, &newTag.Name, &newTag.Slug)

	if err := row.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
		return nil, err
	}

	return &tag, nil
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
