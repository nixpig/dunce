package article

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestArticleRepo(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleRepository){
		// create
		"test create new article":        testCreateNewArticle,
		"test create fails on db errors": testCreateNewArticleFailsOnDbErrors,

		// delete
		"test delete article successfully": testDeleteArticle,
		"test delete article error":        testDeleteArticleError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create database mock")
			}

			data := NewArticleRepository(mock)

			fn(t, mock, data)
		})
	}
}

func testCreateNewArticle(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleRepository) {
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

	newArticle := ArticleNew{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{4},
	}

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

func testCreateNewArticleFailsOnDbErrors(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	mock.ExpectBegin().WillReturnError(errors.New("db_begin_error"))

	article, err := repo.Create(&ArticleNew{
		Title:     "title",
		Subtitle:  "subtitle",
		Slug:      "slug",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    []int{42, 69},
	})

	require.Nil(t, article, "should not return article")
	require.EqualError(t, err, "db_begin_error", "should return db error")

	mock.Reset()
	mock.ExpectationsWereMet()
}

func testDeleteArticle(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	query := `delete from articles_ a using article_tags_ t where a.id_ = t.article_id_ and a.id_ = $1`

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(23).
		WillReturnResult(pgxmock.NewResult("delete", 1))

	err := repo.DeleteById(23)

	mock.Reset()
	mock.ExpectationsWereMet()

	require.Nil(t, err, "should not return error")
}

func testDeleteArticleError(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	query := `delete from articles_ a using article_tags_ t where a.id_ = t.article_id_ and a.id_ = $1`

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(23).
		WillReturnError(errors.New("db_delete_error"))

	err := repo.DeleteById(23)

	mock.Reset()
	mock.ExpectationsWereMet()

	require.EqualError(t, err, "db_delete_error", "should return db error")
}
