package tag

import (
	"context"

	"github.com/nixpig/dunce/internal/pkg/models"
)

type TagData struct {
	db models.Dbconn
}

func NewTagData(db models.Dbconn) TagData {
	return TagData{db}
}

func (u *TagData) Create(tag *TagNew) (*Tag, error) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	var createdTag Tag

	row := u.db.QueryRow(context.Background(), query, tag.Name, tag.Slug)

	if err := row.Scan(&createdTag.Id, &createdTag.Name, &createdTag.Slug); err != nil {
		return nil, err
	}

	return &createdTag, nil
}

func (u *TagData) DeleteById(id int) error {
	query := `delete from tags_ where id_ = $1`

	_, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
