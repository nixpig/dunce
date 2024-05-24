package site

import (
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestSiteRepository(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, repo SiteRepository){
		"test create site key-value": testCreateSiteKeyValue,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db pool")
			}

			repo := NewSitePostgresRepository(mock)

			fn(t, mock, repo)
		})

	}
}

func testCreateSiteKeyValue(t *testing.T, mock pgxmock.PgxPoolIface, repo SiteRepository) {
	query := `insert into site_ (key_, value_) values ($1, $2) returning id_, key_, value_`

	mockRow := mock.
		NewRows([]string{"id_", "key_", "value_"}).
		AddRow(uint(23), "name", "site name")

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("name", "site name").
		WillReturnRows(mockRow)

	got, err := repo.Create("name", "site name")

	require.NoError(t, err, "should not return error")
	require.Equal(t, &SiteKv{
		Id: 23, Key: "name", Value: "site name",
	}, got, "should return created site k/v")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("unmet expectations")
	}
}
