package tag

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
}

type TagData struct {
	db db.Dbconn
}

func NewTagData(db db.Dbconn) TagData {
	return TagData{db}
}

func (u *TagData) create(tag *Tag) (*Tag, error) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var createdTag Tag

	row := u.db.QueryRow(context.Background(), query, tag.Name, tag.Slug)

	if err := row.Scan(&createdTag.Id, &createdTag.Name, &createdTag.Slug); err != nil {
		return nil, err
	}

	return &createdTag, nil
}

func (u *TagData) deleteById(id int) error {
	query := `delete from tags_ where id_ = $1`

	_, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
