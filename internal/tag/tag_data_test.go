package tag_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/tag"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func executeTest(t *testing.T, scenario string, fn func(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData)) {
	t.Run(scenario, func(t *testing.T) {
		db, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create mock db pool")
		}

		defer db.Close()

		data := tag.NewTagData(db)

		fn(t, db, data)
	})
}

func TestTagDataCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData){
		"create new tag": testTagDataCreate,
	}

	for scenario, fn := range scenarios {
		executeTest(t, scenario, fn)
	}
}

func testTagDataCreate(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData) {
	query := `insert into tags_ (name_, slug_) values ($1, $2) returning id_, name_, slug_`

	mockTagRows := mock.NewRows([]string{"id_", "name_", "slug_"}).AddRow(23, "tag_name", "tag_slug")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("tag_name", "tag_slug").WillReturnRows(mockTagRows)

	newTag := tag.NewTag("tag_name", "tag_slug")

	createdTag, err := data.Create(&newTag)

	require.Nil(t, err, "should not error")
	require.Equal(t, &tag.Tag{
		Id:   23,
		Name: "tag_name",
		Slug: "tag_slug",
	}, createdTag, "tag should be saved and match")
}

func TestTagDataDeleteById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData){
		"delete existing tag":     testDeleteExistingTag,
		"delete non-existing tag": testDeleteNonExistingTag,
	}

	for scenario, fn := range scenarios {
		executeTest(t, scenario, fn)
	}
}

func testDeleteExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData) {
	query := `delete from tags_ where id_ = $1`

	mockDeleted := pgxmock.NewResult("delete", 1)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnResult(mockDeleted)

	err := data.DeleteById(23)
	require.Nil(t, err, "should not error")
}

func testDeleteNonExistingTag(t *testing.T, mock pgxmock.PgxPoolIface, data tag.TagData) {
	query := `delete from tags_ where id_ = $1`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnError(errors.New("foo"))

	err := data.DeleteById(23)
	require.Error(t, err)
}
