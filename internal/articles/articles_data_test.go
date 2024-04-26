package articles

import (
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestArticleDataCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleData){
		"test create new article": testCreateNewArticle,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create database mock")
			}

			data := NewArticleData(mock)

			fn(t, mock, data)
		})
	}
}

func testCreateNewArticle(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleData) {
	query := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`

	createdAt := time.Now()
	updatedAt := time.Now()

	mockRow := mock.
		NewRows([]string{"id_", "title_", "subtitle_", "slug_", "body_", "created_at_", "updated_at_"}).
		AddRow(13, "article title", "article subtitle", "article-slug", "Lorem ipsum dolar sit amet...", createdAt, updatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("article title", "article subtitle", "article-slug", "Lorem ipsum dolar sit amet...", createdAt, updatedAt).WillReturnRows(mockRow)

	newArticle := NewArticle(
		"article title",
		"article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet...",
		createdAt,
		updatedAt,
		[]int{},
	)

	createdArticle, err := data.create(&newArticle)

	mock.Reset()
	mock.ExpectationsWereMet()

	require.NoError(t, err, "should not error out")

	require.Equal(t, &Article{
		Id:        13,
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, createdArticle, "should return created article data with id")
}
