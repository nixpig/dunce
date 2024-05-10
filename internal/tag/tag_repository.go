package tag

import (
	"context"
	"errors"

	"github.com/nixpig/dunce/db"
)

type ITagRepository interface {
	Create(tag *TagData) (*Tag, error)
	DeleteById(id int) error
	Exists(tag *TagData) (bool, error)
	GetAll() (*[]Tag, error)
	GetByAttribute(attr, value string) (*Tag, error)
	Update(tag *Tag) (*Tag, error)
}

type TagRepository struct {
	db db.Dbconn
}

func NewTagRepository(db db.Dbconn) TagRepository {
	return TagRepository{
		db: db,
	}
}

func (t TagRepository) Create(tag *TagData) (*Tag, error) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var createdTag Tag

	row := t.db.QueryRow(context.Background(), query, tag.Name, tag.Slug)

	if err := row.Scan(&createdTag.Id, &createdTag.Name, &createdTag.Slug); err != nil {
		return nil, err
	}

	return &createdTag, nil
}

func (t TagRepository) DeleteById(id int) error {
	query := `delete from tags_ where id_ = $1`

	_, err := t.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (t TagRepository) Exists(tag *TagData) (bool, error) {
	checkDuplicatesQuery := `select count(*) from tags_ where slug_ = $1`

	var duplicateCount int

	duplicateRow := t.db.QueryRow(context.Background(), checkDuplicatesQuery, tag.Slug)
	if err := duplicateRow.Scan(&duplicateCount); err != nil {
		return false, err
	}

	if duplicateCount > 0 {
		return true, nil
	}

	return false, nil
}

func (t TagRepository) GetAll() (*[]Tag, error) {
	query := `select id_, name_, slug_ from tags_`

	rows, err := t.db.Query(context.Background(), query)
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

func (t TagRepository) GetByAttribute(attr, value string) (*Tag, error) {
	var query string

	switch attr {
	case "slug":
		query = `select id_, name_, slug_ from tags_ where slug_ = $1`
	default:
		return nil, errors.New("invalid atrribute")
	}

	row := t.db.QueryRow(context.Background(), query, value)

	var tag Tag

	if err := row.Scan(&tag.Id, &tag.Name, &tag.Slug); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (t TagRepository) Update(tag *Tag) (*Tag, error) {
	query := `update tags_ set name_ = $2, slug_ = $3 where id_ = $1 returning id_, name_, slug_`

	row := t.db.QueryRow(context.Background(), query, tag.Id, tag.Name, tag.Slug)

	var updatedTag Tag

	if err := row.Scan(&updatedTag.Id, &updatedTag.Name, &updatedTag.Slug); err != nil {
		return nil, err
	}

	return &updatedTag, nil
}
