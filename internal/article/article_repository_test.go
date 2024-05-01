package article

import (
	"regexp"
	"testing"
	"time"

	"github.com/nixpig/dunce/pkg/logging"
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

			data := NewArticleData(mock, logging.NewLogger())

			fn(t, mock, data)
		})
	}
}

func testCreateNewArticle(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleData) {
	articleInsertQuery := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`
	// tagInsertQuery := `insert into article_tags_ (article_id_, tag_id_) values ($1, $2) returning (tag_id_)`

	createdAt := time.Now()
	updatedAt := time.Now()

	articleMockRow := mock.
		NewRows([]string{"id_", "title_", "subtitle_", "slug_", "body_", "created_at_", "updated_at_"}).
		AddRow(13, "article title", "article subtitle", "article-slug", "Lorem ipsum dolar sit amet...", createdAt, updatedAt)

	// tagMockRow := mock.NewRows([]string{"id_"}).AddRow(69)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(articleInsertQuery)).WithArgs("article title", "article subtitle", "article-slug", "Lorem ipsum dolar sit amet...", createdAt, updatedAt).WillReturnRows(articleMockRow)
	// mock.ExpectQuery(regexp.QuoteMeta(tagInsertQuery)).WithArgs(13, 4).WillReturnRows(tagMockRow)
	mock.ExpectCommit()

	newArticle := NewArticle(
		"article title",
		"article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet...",
		createdAt,
		updatedAt,
		[]int{4},
	)

	createdArticle, err := data.Create(&newArticle)

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
		// TODO: add back once pgxmock supports batch
		// TagIds: []int{4},
	}, createdArticle, "should return created article data with id")
}
