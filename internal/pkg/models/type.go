package models

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mrz1836/go-sanitize"
)

type Type struct {
	Db Dbconn
}

type TypeData struct {
	Id       int    `validate:"required"`
	Name     string `validate:"required,max=255"`
	Template string `validate:"required,max=255"`
	Slug     string `validate:"required,slug,max=255"`
}

type NewTypeData struct {
	Name     string `validate:"required,max=255"`
	Template string `validate:"required,max=255"`
	Slug     string `validate:"required,slug,max=255"`
}

type UpdateTypeData struct {
	Id       int    `validate:"required"`
	Name     string `validate:"required,max=255"`
	Template string `validate:"required,max=255"`
	Slug     string `validate:"required,slug,max=255"`
}

func (t *Type) GetAll() (*[]TypeData, error) {
	query := `select id_, name_, template_, slug_ from types_`

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
	sanitisedTypeData := NewTypeData{
		Name:     sanitize.XSS(newType.Name),
		Template: sanitize.URI(newType.Template),
		Slug:     sanitize.PathName(newType.Slug),
	}
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedTypeData); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	checkDuplicatesQuery := `select count(*) from types_ where name_ = $1 or template_ = $2 or slug_ = $3`

	var duplicateCount int

	duplicateRow := t.Db.QueryRow(context.Background(), checkDuplicatesQuery, &sanitisedTypeData.Name, &sanitisedTypeData.Template, &sanitisedTypeData.Slug)

	if err := duplicateRow.Scan(&duplicateCount); err != nil {
		return nil, err
	}

	if duplicateCount > 0 {
		return nil, fmt.Errorf("duplicate type: '%s' '%s' '%s'", sanitisedTypeData.Name, sanitisedTypeData.Template, sanitisedTypeData.Slug)
	}

	query := `insert into types_ (name_, template_, slug_) values ($1, $2, $3) returning id_, name_, template_, slug_`

	row := t.Db.QueryRow(context.Background(), query, &sanitisedTypeData.Name, &sanitisedTypeData.Template, &sanitisedTypeData.Slug)

	var createdType TypeData

	if err := row.Scan(&createdType.Id, &createdType.Name, &createdType.Template, &createdType.Slug); err != nil {
		return nil, err
	}

	return &createdType, nil
}

func (t *Type) GetById(id int) (*TypeData, error) {
	query := `select id_, name_, template_, slug_ from types_ where id_ = $1`

	row := t.Db.QueryRow(context.Background(), query, id)

	var typeData TypeData

	if err := row.Scan(&typeData.Id, &typeData.Name, &typeData.Template, &typeData.Slug); err != nil {
		return nil, err
	}

	return &typeData, nil
}

func (t *Type) DeleteById(id int) error {
	query := `delete from types_ where id_ = $1`

	res, err := t.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("nothing deleted")
	}

	return nil
}

func (t *Type) UpdateById(id int, typeData UpdateTypeData) (*TypeData, error) {
	sanitisedTypeData := NewTypeData{
		Name:     sanitize.XSS(typeData.Name),
		Template: sanitize.URI(typeData.Template),
		Slug:     sanitize.PathName(typeData.Slug),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", ValidateSlug)

	if err := validate.Struct(sanitisedTypeData); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	checkDuplicatesQuery := `select count(*) from types_ where (name_ = $2 or slug_ = $3) and id_ <> $1`

	var duplicateCount int

	res := t.Db.QueryRow(context.Background(), checkDuplicatesQuery, id, &sanitisedTypeData.Name, &sanitisedTypeData.Slug)

	if err := res.Scan(&duplicateCount); err != nil {
		return nil, err
	}

	if duplicateCount > 0 {
		return nil, fmt.Errorf("duplicate type: '%s' '%s' '%s'", sanitisedTypeData.Name, sanitisedTypeData.Template, sanitisedTypeData.Slug)
	}

	query := `update types_ set name_ = $2, template_ = $3, slug_ = $4 where id_ = $1 returning id_, name_, template_, slug_`

	updated := t.Db.QueryRow(context.Background(), query, id, &sanitisedTypeData.Name, &sanitisedTypeData.Template, &sanitisedTypeData.Slug)

	var updatedType TypeData

	if err := updated.Scan(&updatedType.Id, &updatedType.Name, &updatedType.Template, &updatedType.Slug); err != nil {
		return nil, err
	}

	return &updatedType, nil
}
