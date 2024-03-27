package models_test

import (
	"fmt"
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
		"fail to create duplicate":                       testFailCreateDuplicateTag,
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
		WithArgs(newTag.Name, newTag.Slug).
		WillReturnRows(noDuplicateRows)

	mock.
		ExpectQuery(regexp.QuoteMeta(insertQuery)).
		WithArgs(newTag.Name, newTag.Slug).
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
	newTag := models.TagData{
		Name: "<script>alert('xss');</script>",
		Slug: "\\ tag-& slug*($$)",
	}

	sanitisedTagData := models.TagData{
		Name: ">alert('xss');</",
		Slug: "tag-slug",
	}

	duplicateQuery := "select count(*) from tags_ where name_ = $1 or slug_ = $2"
	insertQuery := "insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_"

	noDuplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(0)

	insertedRow := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(1, ">alert('xss');</", "tag-slug")

	mock.
		ExpectQuery(regexp.QuoteMeta(duplicateQuery)).
		WithArgs(sanitisedTagData.Name, sanitisedTagData.Slug).
		WillReturnRows(noDuplicateRows)

	mock.
		ExpectQuery(regexp.QuoteMeta(insertQuery)).
		WithArgs(sanitisedTagData.Name, sanitisedTagData.Slug).
		WillReturnRows(insertedRow)

	tag, err := models.Query.Tag.Create(newTag)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Tag{
		Id: 1,
		TagData: models.TagData{
			Name: ">alert('xss');</",
			Slug: "tag-slug",
		},
	}, tag, "should return created sanitised tag")
}

func testFailCreateDuplicateTag(t *testing.T, mock pgxmock.PgxPoolIface) {
	newTag := models.TagData{
		Name: "newtag",
		Slug: "new-tag",
	}

	existingTagRows := mock.NewRows([]string{"count"}).AddRow(1)

	duplicateTagQuery := "select count(*) from tags_ where name_ = $1 or slug_ = $2"

	mock.ExpectQuery(regexp.QuoteMeta(duplicateTagQuery)).WithArgs(newTag.Name, newTag.Slug).WillReturnRows(existingTagRows)

	tag, err := models.Query.Tag.Create(newTag)
	require.Nil(t, tag, "should not return a tag")
	require.Equal(t, fmt.Errorf("duplicate tag: '%s' '%s'", newTag.Name, newTag.Slug), err, "should return an error")
}

func testFailCreateTagInvalidSlug(t *testing.T, mock pgxmock.PgxPoolIface) {
	newTag := models.TagData{
		Name: "newtag",
		Slug: "long-slug-long-slug-long-slug-long-slug-long-slug-long-slug-long-slug-long-slug-long-slug-long-slug-long-slug",
	}

	tag, err := models.Query.Tag.Create(newTag)
	require.Nil(t, tag, "should not return a tag")
	require.NotNil(t, err, "should return an error")
}

func testFailCreateTagInvalidName(t *testing.T, mock pgxmock.PgxPoolIface) {
	newTag := models.TagData{
		Name: "longnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongnamelongname",
		Slug: "long-slug",
	}

	tag, err := models.Query.Tag.Create(newTag)
	require.Nil(t, tag, "should not return a tag")
	require.NotNil(t, err, "should return an error")
}
