package models_test

import (
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreateTag(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"successfully create new tag":                    testCreateNewTag,
		"successfully sanitise name and slug of new tag": testSanitiseNewTag,
		"fail to create tag with duplicate name":         testFailCreateTagDuplicateName,
		"fail to create tag with duplicate slug":         testFailCreateTagDuplicateSlug,
		"fail to create tag with invalid name":           testFailCreateTagInvalidName,
		"fail to create tag with invalid slug":           testFailCreateTagInvalidSlug,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("failed to create mock database pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func testCreateNewTag(t *testing.T, mock pgxmock.PgxPoolIface) {
	newTag := models.TagData{
		Name: "tagname",
		Slug: "tag-slug",
	}

	duplicateQuery := "select count(*) from tags_ where name_ = $1 or slug_ = $2"
	insertQuery := "insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_"

	noDuplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(0)

	insertedRow := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(1, "tagname", "tag-slug")

	mock.
		ExpectQuery(regexp.QuoteMeta(duplicateQuery)).
		WithArgs(&newTag.Name, &newTag.Slug).
		WillReturnRows(noDuplicateRows)

	mock.
		ExpectQuery(regexp.QuoteMeta(insertQuery)).
		WithArgs(&newTag.Name, &newTag.Slug).
		WillReturnRows(insertedRow)

	tag, err := models.Query.Tag.Create(newTag)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Tag{
		Id: 1,
		TagData: models.TagData{
			Name: "tagname",
			Slug: "tag-slug",
		},
	}, tag, "should return created tag")
}

func testSanitiseNewTag(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailCreateTagDuplicateName(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailCreateTagDuplicateSlug(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailCreateTagInvalidSlug(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailCreateTagInvalidName(t *testing.T, mock pgxmock.PgxPoolIface) {

}
