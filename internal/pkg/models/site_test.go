package models_test

import (
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestSiteModel(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get site data":            testGetSiteData,
		"update valid site data":   testUpdateValidSiteData,
		"update invalid site data": testUpdateInvalidSiteData,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func testGetSiteData(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select name_, description_, url_, owner_ from site_ limit 1`

	mockResult := mock.
		NewRows([]string{"name_", "description_", "url_", "owner_"}).
		AddRow("Test", "Test description", "https://t.com/f", 23)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockResult)

	siteData, err := models.Query.Site.Get()
	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Site{
		Name:        "Test",
		Description: "Test description",
		Url:         "https://t.com/f",
		Owner:       23,
	}, siteData, "should return site data")
}

func testUpdateValidSiteData(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `update site_ set name_ = $1, description_ = $2, url_ = $3, owner_ = $4 returning name_, description_, url_, owner_`

	siteUpdate := models.Site{
		Name:        "update name",
		Description: "update description",
		Url:         "https://t.com/f",
		Owner:       23,
	}

	result := mock.
		NewRows([]string{"name_", "description_", "url_", "owner_"}).
		AddRow("update name", "update description", "https://t.com/f", 23)

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(
			&siteUpdate.Name,
			&siteUpdate.Description,
			&siteUpdate.Url,
			&siteUpdate.Owner,
		).
		WillReturnRows(result)

	updated, err := models.Query.Site.Update(siteUpdate)
	require.Nil(t, err, "should not error")
	require.Equal(t, &models.Site{
		Name:        "update name",
		Description: "update description",
		Url:         "https://t.com/f",
		Owner:       23,
	}, updated, "should return updated site details")
}

func testUpdateInvalidSiteData(t *testing.T, mock pgxmock.PgxPoolIface) {
	siteUpdate := models.Site{
		Name:        "long name long name long name long name long name long name ",
		Description: "long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description long description ",
		Url:         "not-a-url",
		Owner:       23,
	}

	updated, err := models.Query.Site.Update(siteUpdate)
	require.Nil(t, updated, "should not update site details")
	require.NotNil(t, err, "should return error")
}
