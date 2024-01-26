package models

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Type struct {
	Db Dbconn
}

type TypeData struct {
	Id       int    `validate:"required"`
	Name     string `validate:"required,max=255"`
	Template string `validate:"required,max=255"`
	Slug     string `validate:"required,max=255"`
}

type NewTypeData struct {
	Name     string `validate:"required,max=255"`
	Template string `validate:"required,max=255"`
	Slug     string `validate:"required,max=255"`
}

func (t *Type) GetAll() (*[]TypeData, error) {
	query := `select id_, name_, template_, slug_ from type_`

	rows, err := t.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var types []TypeData

	for rows.Next() {
		var typeData TypeData

		if err := rows.Scan(&typeData.Id, &typeData.Name, &typeData.Template, &typeData.Slug); err != nil {
			return nil, err
		}

		types = append(types, typeData)
	}

	return &types, nil
}

func (t *Type) Create(newType NewTypeData) (*TypeData, error) {
	validate := validator.New()

	if err := validate.Struct(newType); err != nil {
		return nil, err
	}

	query := `insert into type_ (name_, template_, slug_) values ($1, $2, $3) returning id_, name_, template_, slug_`

	row := t.Db.QueryRow(context.Background(), query, &newType.Name, &newType.Template, &newType.Slug)

	var createdType TypeData

	if err := row.Scan(&createdType.Id, &createdType.Name, &createdType.Template, &createdType.Slug); err != nil {
		return nil, err
	}

	return &createdType, nil
}
