package tags

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestTagDataUpdate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"update tag": testTagDataUpdateTag,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			db, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("failed to create mock db pool")
			}

			defer db.Close()

			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func TestTagDataGetBySlug(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"get existing tag by slug":     testGetExistingTagBySlug,
		"get non-existent tag by slug": testGetNonExistentTagBySlug,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			db, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("failed to create mock db")
			}

			defer db.Close()

			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func TestTagDataGetAll(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"get all (multiple results)": testGetAllTagsMultipleResults,
		"get all (no results)":       testGetAllTagsNoResults,
		"get all (single result)":    testGetAllTagsSingleResult,
	}

	for scenario, fn := range scenarios {
		db, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create mock db pool")
		}

		defer db.Close()

		t.Run(scenario, func(t *testing.T) {
			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func TestTagDataCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"create new tag": testCreateValidTag,
	}

	for scenario, fn := range scenarios {
		db, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create mock db pool")
		}

		defer db.Close()

		t.Run(scenario, func(t *testing.T) {
			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func TestTagDataDeleteById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"delete existing tag":     testDeleteExistingTag,
		"delete non-existing tag": testDeleteNonExistingTag,
	}

	for scenario, fn := range scenarios {
		db, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create mock db pool")
		}

		defer db.Close()

		t.Run(scenario, func(t *testing.T) {
			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func TestTagDataExists(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data TagData){
		"check existing tag exists":     testTagExists,
		"check existing tag not exists": testTagNotExists,
		"check existing tag error":      testTagExistsError,
	}

	for scenario, fn := range scenarios {
		db, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create mock db pool")
		}

		defer db.Close()

		t.Run(scenario, func(t *testing.T) {
			data := NewTagData(db)

			fn(t, db, data)
		})
	}
}

func testCreateValidTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	mockTagRows := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRow(23, "tag_name", "tag_slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("tag_name", "tag_slug").WillReturnRows(mockTagRows)

	newTag := NewTag("tag_name", "tag_slug")

	createdTag, err := data.create(&newTag)

	require.NoError(t, err, "should not error")
	require.Equal(t, &Tag{
		Id:   23,
		Name: "tag_name",
		Slug: "tag_slug",
	}, createdTag, "tag should be saved and match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testCreateInvalidTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("database_error"))

	newTag := NewTag(
		"some_really_long_namesome_really_long_namesome_really_long_name",
		"some-really-long-tagsome-really-long-tagsome-really-long-tagsome-really-long-tag",
	)

	createdTag, err := data.create(&newTag)
	require.Nil(t, createdTag, "should not create invalid tag")
	require.EqualError(t, err, "database_error", "should return the error from database")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testDeleteExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `delete from tags_ where id_ = $1`

	mockDeleted := pgxmock.NewResult("delete", 1)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnResult(mockDeleted)

	err := data.deleteById(23)
	require.NoError(t, err, "should not error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testDeleteNonExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `delete from tags_ where id_ = $1`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnError(errors.New("database_error"))

	err := data.deleteById(23)
	require.EqualError(t, err, "database_error", "should return error from database")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testTagExists(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	duplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnRows(duplicateRows)

	exists, err := data.exists(&mockDuplicateTag)

	require.NoError(t, err, "should not return error")
	require.True(t, exists, "should return true")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testTagNotExists(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	duplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnRows(duplicateRows)

	exists, err := data.exists(&mockDuplicateTag)

	require.NoError(t, err, "should not return error")
	require.False(t, exists, "should return false")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testTagExistsError(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnError(errors.New("database_error"))

	_, err := data.exists(&mockDuplicateTag)

	require.EqualError(t, err, "database_error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testGetAllTagsNoResults(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select id_, name_, slug_ from tags_`

	mockEmptyRows := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRows()

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockEmptyRows)

	tags, err := data.GetAll()

	require.NoError(t, err, "should not return error")
	require.Empty(t, tags, "should return zero results")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testGetAllTagsMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select id_, name_, slug_ from tags_`

	singleResult := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(23, "tagname1", "tag-slug-1").
		AddRow(42, "tagname2", "tag-slug-2").
		AddRow(69, "tagname3", "tag-slug-3")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(singleResult)

	tags, err := data.GetAll()
	require.Equal(t, &[]Tag{
		{
			Id:   23,
			Name: "tagname1",
			Slug: "tag-slug-1",
		},
		{
			Id:   42,
			Name: "tagname2",
			Slug: "tag-slug-2",
		},
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, tags, "should return all tag results")
	require.NoError(t, err, "should not return an error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testGetAllTagsSingleResult(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select id_, name_, slug_ from tags_`

	singleResult := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(23, "tagname", "tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(singleResult)

	tags, err := data.GetAll()
	require.Equal(t, &[]Tag{
		{

			Id:   23,
			Name: "tagname",
			Slug: "tag-slug",
		},
	}, tags, "should return tag result")
	require.NoError(t, err, "should not return an error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testGetExistingTagBySlug(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select id_, name_, slug_ from tags_ where slug_ = $1`

	mockRow := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRow(23, "tagname", "tag-slug")

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("tag-slug").
		WillReturnRows(mockRow)

	tag, err := data.getBySlug("tag-slug")

	require.NoError(t, err, "should not return error")
	require.Equal(t, &Tag{
		Id:   23,
		Name: "tagname",
		Slug: "tag-slug",
	}, tag, "should return tag")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testGetNonExistentTagBySlug(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select id_, name_, slug_ from tags_ where slug_ = $1`

	mockRow := mock.NewRows([]string{"id_", "name_", "slug_"})

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("tag-slug").
		WillReturnRows(mockRow)

	tag, err := data.getBySlug("tag-slug")

	require.Error(t, err, "should return error")
	require.Nil(t, tag, "should not return any tag")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met")
	}
}

func testTagDataUpdateTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `update tags_ set name_ = $2, slug_ = $3 where id_ = $1 returning id_, name_, slug_`

	mockRes := mock.NewRows([]string{"id_", "name_", "id_"}).AddRow(23, "tagname", "tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(23, "tagname", "tag-slug").WillReturnRows(mockRes)

	tagUpdate := NewTagWithId(23, "tagname", "tag-slug")

	tag, err := data.update(&tagUpdate)

	require.NoError(t, err, "should not error")
	require.Equal(t, &Tag{
		Id:   23,
		Name: "tagname",
		Slug: "tag-slug",
	}, tag, "should return updated tag")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("mock expectations were not met")
	}
}
