package tag

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

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

	require.Nil(t, err, "should not error")
	require.Equal(t, &Tag{
		Id:   23,
		Name: "tag_name",
		Slug: "tag_slug",
	}, createdTag, "tag should be saved and match")
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

func testDeleteExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `delete from tags_ where id_ = $1`

	mockDeleted := pgxmock.NewResult("delete", 1)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnResult(mockDeleted)

	err := data.deleteById(23)
	require.Nil(t, err, "should not error")
}

func testDeleteNonExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `delete from tags_ where id_ = $1`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnError(errors.New("database_error"))

	err := data.deleteById(23)
	require.EqualError(t, err, "database_error", "should return error from database")
}

func testTagExists(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	duplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnRows(duplicateRows)

	exists, err := data.exists(&mockDuplicateTag)

	require.Nil(t, err, "should not return error")
	require.True(t, exists, "should return true")
}

func testTagNotExists(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	duplicateRows := mock.
		NewRows([]string{"count"}).
		AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnRows(duplicateRows)

	exists, err := data.exists(&mockDuplicateTag)

	require.Nil(t, err, "should not return error")
	require.False(t, exists, "should return false")
}

func testTagExistsError(t *testing.T, mock pgxmock.PgxPoolIface, data TagData) {
	query := `select count(*) from tags_ where name_ = $1 or slug_ = $2`

	mockDuplicateTag := NewTagWithId(23, "existing tag name", "existing-tag-slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockDuplicateTag.Name, mockDuplicateTag.Slug).WillReturnError(errors.New("database_error"))

	_, err := data.exists(&mockDuplicateTag)

	require.EqualError(t, err, "database_error")
}
