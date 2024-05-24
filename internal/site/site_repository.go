package site

import (
	"context"

	"github.com/nixpig/dunce/db"
)

type SiteRepository interface {
	Create(key, value string) (*SiteKv, error)
}

type sitePostgresRepository struct {
	db db.Dbconn
}

func NewSitePostgresRepository(db db.Dbconn) sitePostgresRepository {
	return sitePostgresRepository{db}
}

func (s sitePostgresRepository) Create(key, value string) (*SiteKv, error) {
	query := `insert into site_ (key_, value_) values ($1, $2) returning key_, value_`

	row := s.db.QueryRow(context.Background(), query, key, value)

	var created SiteKv

	if err := row.Scan(&created.Key, &created.Value); err != nil {
		return nil, err
	}

	return &created, nil
}
