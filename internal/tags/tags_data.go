package tags

import (
	"context"

	"github.com/nixpig/dunce/db"
)

type Tag struct {
	Id   int
	Name string `validate:"required,min=5,max=30"`
	Slug string `validate:"required,slug,min=5,max=50"`
}

func NewTag(name, slug string) Tag {
	return Tag{Name: name, Slug: slug}
}

func NewTagWithId(id int, name, slug string) Tag {
	return Tag{Id: id, Name: name, Slug: slug}
}

type TagDataInterface interface {
	create(tag *Tag) (*Tag, error)
	deleteById(id int) error
	exists(tag *Tag) (bool, error)
	getAll() (*[]Tag, error)
}

type TagData struct {
	db db.Dbconn
}

func NewTagData(db db.Dbconn) TagData {
	return TagData{db}
}

func (t *TagData) create(tag *Tag) (*Tag, error) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var createdTag Tag

	row := t.db.QueryRow(context.Background(), query, tag.Name, tag.Slug)

	if err := row.Scan(&createdTag.Id, &createdTag.Name, &createdTag.Slug); err != nil {
		return nil, err
	}

	return &createdTag, nil
}

func (t *TagData) deleteById(id int) error {
	query := `delete from tags_ where id_ = $1`

	_, err := t.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TagData) exists(tag *Tag) (bool, error) {
	checkDuplicatesQuery := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	var duplicateCount int

	duplicateRow := t.db.QueryRow(context.Background(), checkDuplicatesQuery, tag.Name, tag.Slug)
	if err := duplicateRow.Scan(&duplicateCount); err != nil {
		return false, err
	}

	if duplicateCount > 0 {
		return true, nil
	}

	return false, nil
}

func (t *TagData) getAll() (*[]Tag, error) {
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
