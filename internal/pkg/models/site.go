package models

import (
	"context"
)

type SiteData struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=255"`
	Url         string `json:"url" validate:"required,url,max=255"`
	Owner       int    `json:"owner" validate:"required"`
}

type Site struct {
	Db Dbconn
}

func (s *Site) Get() (*SiteData, error) {
	query := `select name_, description_, url_, owner_ from site_ limit 1`

	row := s.Db.QueryRow(context.Background(), query)

	var siteData SiteData

	if err := row.Scan(&siteData.Name, &siteData.Description, &siteData.Url, &siteData.Owner); err != nil {
		return nil, err
	}

	return &siteData, nil
}

func (s *Site) Update(data SiteData) (*SiteData, error) {
	query := `update site_ set name_ = $1, description_ = $2, url_ = $3, owner_ = $4 returning name_, description_, url_, owner_`
	row := s.Db.QueryRow(context.Background(), query, data.Name, data.Description, data.Url, data.Owner)

	var updatedData SiteData

	if err := row.Scan(&updatedData.Name, &updatedData.Description, &updatedData.Url, &updatedData.Owner); err != nil {
		return nil, err
	}

	return &updatedData, nil
}
