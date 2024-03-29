package models_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestUpdateTag(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"successfully update tag":                testUpdateTag,
		"sanitise name and slug of updated tag":  testSanitiseUpdatedTag,
		"fail to update tag with duplicate name": testFailUpdateTagDuplicateName,
		"fail to update tag with duplicate slug": testFailUpdateTagDuplicateSlug,
		"fail to update tag with invalid name":   testFailUpdateTagInvalidName,
		"fail to update tag with invalid slug":   testFailUpdateTagInvalidSlug,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("failed to create mock pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetTagBySlug(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get tag when slug exists":         testGetExistingTagBySlug,
		"get tag when slug does not exist": testGetNonexistentTagBySlug,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("failed to create mock pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetTagById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get tag by id when exists":         testGetExistingTagById,
		"get tag by id when does not exist": testGetNonexistentTagById,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetAllTags(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get all tags when none exist":       testGetAllTagsNoResults,
		"get all tags when one exists":       testGetAllTagsSingleResult,
		"get all tags when multiple results": testGetAllTagsMultipleResults,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("failed to create mock connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

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

func testGetAllTagsNoResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_`

	emptyResults := mock.NewRows([]string{"id_", "name_", "slug_"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(emptyResults)

	tags, err := models.Query.Tag.GetAll()
	require.Empty(t, tags, "should return empty result set")
	require.Nil(t, err, "should not return an error")
}

func testGetAllTagsSingleResult(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_`

	singleResult := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(23, "tagname", "tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(singleResult)

	tags, err := models.Query.Tag.GetAll()
	require.Equal(t, &[]models.Tag{
		{
			Id: 23,
			TagData: models.TagData{
				Name: "tagname",
				Slug: "tag-slug",
			},
		},
	}, tags, "should return tag result")
	require.Nil(t, err, "should not return an error")
}

func testGetAllTagsMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_`

	singleResult := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(23, "tagname1", "tag-slug-1").
		AddRow(42, "tagname2", "tag-slug-2").
		AddRow(69, "tagname3", "tag-slug-3")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(singleResult)

	tags, err := models.Query.Tag.GetAll()
	require.Equal(t, &[]models.Tag{
		{
			Id: 23,
			TagData: models.TagData{
				Name: "tagname1",
				Slug: "tag-slug-1",
			},
		},
		{
			Id: 42,
			TagData: models.TagData{
				Name: "tagname2",
				Slug: "tag-slug-2",
			},
		},
		{
			Id: 69,
			TagData: models.TagData{
				Name: "tagname3",
				Slug: "tag-slug-3",
			},
		},
	}, tags, "should return all tag results")
	require.Nil(t, err, "should not return an error")
}

func testGetExistingTagById(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_ where id_ = $1`

	row := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRow(23, "tagname", "tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(23).WillReturnRows(row)

	tag, err := models.Query.Tag.GetById(23)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Tag{
		Id: 23,
		TagData: models.TagData{
			Name: "tagname",
			Slug: "tag-slug",
		},
	}, tag, "should return corresponding tag")
}

func testGetNonexistentTagById(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_ where id_ = $1`

	noResults := mock.NewRows([]string{"id_", "name_", "slug_"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(23).WillReturnRows(noResults)

	tag, err := models.Query.Tag.GetById(23)

	require.Nil(t, tag, "should not return a tag")
	require.NotNil(t, err, "should return an error")
}

func testGetExistingTagBySlug(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_ where slug_ = $1`

	row := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRow(23, "tagname", "tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("tag-slug").WillReturnRows(row)

	tag, err := models.Query.Tag.GetBySlug("tag-slug")
	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Tag{
		Id: 23,
		TagData: models.TagData{
			Name: "tagname",
			Slug: "tag-slug",
		},
	}, tag, "should return found tag")
}

func testGetNonexistentTagBySlug(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `select id_, name_, slug_ from tags_ where slug_ = $1`

	emptyRow := mock.NewRows([]string{"id_", "name_", "slug_"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("foo").WillReturnRows(emptyRow)

	tag, err := models.Query.Tag.GetBySlug("foo")

	require.Nil(t, tag, "should not return tag")
	require.NotNil(t, err, "should return error")
}

func testUpdateTag(t *testing.T, mock pgxmock.PgxPoolIface) {
	tagUpdate := models.TagData{
		Name: "updatetag",
		Slug: "update-tag",
	}

	duplicateQuery := `select count(*) from tags_ where (name_ = $2 or slug_ = $3) and id_ <> $1`
	updateQuery := `update tags_ set name_ = $2, slug_ = $3 where id_ = $1 returning id_, name_, slug_`

	zeroCount := mock.NewRows([]string{"count"}).AddRow(0)

	updatedRow := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(23, "updatetag", "update-tag")

	mock.
		ExpectQuery(regexp.QuoteMeta(duplicateQuery)).
		WithArgs(23, &tagUpdate.Name, &tagUpdate.Slug).
		WillReturnRows(zeroCount)

	mock.
		ExpectQuery(regexp.QuoteMeta(updateQuery)).
		WithArgs(23, &tagUpdate.Name, &tagUpdate.Slug).
		WillReturnRows(updatedRow)

	tag, err := models.Query.Tag.UpdateById(23, tagUpdate)
	require.Nil(t, err, "should not return error")
	require.Equal(t, &models.Tag{
		Id: 23,
		TagData: models.TagData{
			Name: "updatetag",
			Slug: "update-tag",
		},
	}, tag, "should return updated tag")
}

func testSanitiseUpdatedTag(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailUpdateTagDuplicateName(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailUpdateTagDuplicateSlug(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailUpdateTagInvalidName(t *testing.T, mock pgxmock.PgxPoolIface) {

}

func testFailUpdateTagInvalidSlug(t *testing.T, mock pgxmock.PgxPoolIface) {

}
