package article

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/nixpig/dunce/internal/tag"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestArticleRepo(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository){
		"test create article (success)":                            testArticleRepoCreateNewArticle,
		"test create article (handle db error)":                    testArticleRepoCreateNewArticleFailsOnDbErrors,
		"test delete article (success)":                            testArticleRepoDeleteArticle,
		"test delete article (handle db error)":                    testArticleRepoDeleteArticleError,
		"test get article by slug (success)":                       testArticleRepoGetArticleBySlug,
		"test get article (error - non-implemented attr)":          testArticleRepoGetArticleByInvalidAttr,
		"test get article (error - article db error)":              testArticleRepoGetArticleByAttrArticleDbError,
		"test get article (error - tags db error)":                 testArticleRepoGetArticleByAttrTagsDbError,
		"test get article (error - tags scan error)":               testArticleRepoGetArticleByAttrTagsScanError,
		"test get many (error - by unknown attribute)":             testArticleRepoGetManyByUnknownAttr,
		"test get many (success - single result - by tag slug)":    testArticleRepoGetManyArticlesByTagSlugSingleResult,
		"test get many (success - multiple results - by tag slug)": testArticleRepoGetManyArticlesByTagSlugMultipleResults,
		"test get all (success - single result)":                   testArticleRepoGetAllArticlesSingleResult,
		"test get all (success - multiple results)":                testArticleRepoGetAllArticlesMultipleResults,

		// read
		// update
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create database mock")
			}

			data := NewArticlePostgresRepository(mock)

			fn(t, mock, data)
		})
	}
}

func testArticleRepoCreateNewArticle(t *testing.T, mock pgxmock.PgxPoolIface, data ArticleRepository) {
	articleInsertQuery := `insert into articles_ (title_, subtitle_, slug_, body_, created_at_, updated_at_) values ($1, $2, $3, $4, $5, $6) returning id_, title_, subtitle_, slug_, body_, created_at_, updated_at_`

	createdAt := time.Now()
	updatedAt := time.Now()

	articleMockRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
	}).AddRow(
		13,
		"article title",
		"article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet...",
		createdAt,
		updatedAt,
	)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(articleInsertQuery)).WithArgs("article title", "article subtitle", "article-slug", "Lorem ipsum dolar sit amet...", createdAt, updatedAt).WillReturnRows(articleMockRow)
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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met: ", err)
	}

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

func testArticleRepoCreateNewArticleFailsOnDbErrors(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met: ", err)
	}
}

func testArticleRepoDeleteArticle(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	query := `delete from articles_ a where a.id_ = $1`

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(23).
		WillReturnResult(pgxmock.NewResult("delete", 1))

	err := repo.DeleteById(23)

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met: ", err)
	}

	require.Nil(t, err, "should not return error")
}

func testArticleRepoDeleteArticleError(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	query := `delete from articles_ a where a.id_ = $1`

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(23).
		WillReturnError(errors.New("db_delete_error"))

	err := repo.DeleteById(23)

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met: ", err)
	}

	require.EqualError(t, err, "db_delete_error", "should return db error")
}

func testArticleRepoGetArticleBySlug(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`

	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour * 24)

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
		"tag_ids_",
	}).AddRow(
		23,
		"Article title",
		"Article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet",
		createdAt,
		updatedAt,
		"42,69",
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("tag-slug").
		WillReturnRows(mockArticleRow)

	tagQuery := `select id_, name_, slug_ from tags_ where id_ = 42 or id_ = 69`

	mockTagRows := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(42, "tag one", "tag-one").
		AddRow(69, "tag two", "tag-two")

	mock.ExpectQuery(regexp.QuoteMeta(tagQuery)).WillReturnRows(mockTagRows)

	got, err := repo.GetByAttribute("slug", "tag-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, err, "should not return error")

	require.Equal(t, &Article{
		Id:        23,
		Title:     "Article title",
		Subtitle:  "Article subtitle",
		Slug:      "article-slug",
		Body:      "Lorem ipsum dolar sit amet",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags: []tag.Tag{
			{Id: 42, Name: "tag one", Slug: "tag-one"},
			{Id: 69, Name: "tag two", Slug: "tag-two"},
		},
	}, got, "should return article")
}

func testArticleRepoGetArticleByInvalidAttr(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	got, err := repo.GetByAttribute("foo", "bar")

	require.Nil(t, got, "should not return any article")
	require.EqualError(t, err, "invalid attribute", "should return invalid attribute error")
}

func testArticleRepoGetArticleByAttrArticleDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("tag-slug").
		WillReturnError(errors.New("article_db_error"))

	got, err := repo.GetByAttribute("slug", "tag-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, got, "should not return article")
	require.EqualError(t, err, "article_db_error", "should return db error")
}

func testArticleRepoGetArticleByAttrTagsDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`

	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour * 24)

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
		"tag_ids_",
	}).AddRow(
		23,
		"Article title",
		"Article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet",
		createdAt,
		updatedAt,
		"42,69",
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("tag-slug").
		WillReturnRows(mockArticleRow)

	tagQuery := `select id_, name_, slug_ from tags_ where id_ = 42 or id_ = 69`

	mock.
		ExpectQuery(regexp.QuoteMeta(tagQuery)).
		WillReturnError(errors.New("tags_db_error"))

	got, err := repo.GetByAttribute("slug", "tag-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, got, "should not return tags")

	require.EqualError(t, err, "tags_db_error", "should return db error")
}

func testArticleRepoGetManyArticlesByTagSlugSingleResult(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_ from articles_ a inner join article_tags_ at on a.id_ = at.article_id_ inner join tags_ t on at.tag_id_ = t.id_ where t.slug_ = $1`

	createdAt := time.Now().Add(time.Hour * -12)
	updatedAt := time.Now()

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
	}).AddRow(
		23,
		"Article title",
		"Article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet",
		createdAt,
		updatedAt,
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("some-slug").
		WillReturnRows(mockArticleRow)

	got, err := repo.GetManyByAttribute("tagSlug", "some-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Article{{
		Id:        23,
		Title:     "Article title",
		Subtitle:  "Article subtitle",
		Slug:      "article-slug",
		Body:      "Lorem ipsum dolar sit amet",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags:      []tag.Tag(nil),
	}}, got, "should return the article")
}

func testArticleRepoGetManyByUnknownAttr(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	got, err := repo.GetManyByAttribute("foo", "bar")

	require.Nil(t, got, "shouldn't return any articles")
	require.EqualError(t, err, "unsupported attribute", "should return unsupported attribute error")
}

func testArticleRepoGetManyArticlesByTagSlugMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_ from articles_ a inner join article_tags_ at on a.id_ = at.article_id_ inner join tags_ t on at.tag_id_ = t.id_ where t.slug_ = $1`

	createdAt := time.Now().Add(time.Hour * -12)
	updatedAt := time.Now()

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
	}).AddRow(
		23,
		"Article title one",
		"Article subtitle one",
		"article-slug-one",
		"Lorem ipsum dolar sit amet one",
		createdAt,
		updatedAt,
	).AddRow(
		42,
		"Article title two",
		"Article subtitle two",
		"article-slug-two",
		"Lorem ipsum dolar sit amet two",
		createdAt,
		updatedAt,
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("some-slug").
		WillReturnRows(mockArticleRow)

	got, err := repo.GetManyByAttribute("tagSlug", "some-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Article{
		{
			Id:        23,
			Title:     "Article title one",
			Subtitle:  "Article subtitle one",
			Slug:      "article-slug-one",
			Body:      "Lorem ipsum dolar sit amet one",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags:      []tag.Tag(nil),
		},
		{
			Id:        42,
			Title:     "Article title two",
			Subtitle:  "Article subtitle two",
			Slug:      "article-slug-two",
			Body:      "Lorem ipsum dolar sit amet two",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags:      []tag.Tag(nil),
		},
	}, got, "should return the article")
}

func testArticleRepoGetAllArticlesSingleResult(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	tagQuery := `select id_, name_, slug_ from tags_`

	mockTagRows := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(42, "tag one", "tag-one").
		AddRow(69, "tag two", "tag-two")

	mock.ExpectQuery(regexp.QuoteMeta(tagQuery)).WillReturnRows(mockTagRows)

	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour * 24)

	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
		"tag_ids_",
	}).AddRow(
		23,
		"Article title",
		"Article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet",
		createdAt,
		updatedAt,
		"42,69",
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WillReturnRows(mockArticleRow)

	got, err := repo.GetAll()

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, err, "should not return error")

	require.Equal(t, &[]Article{
		{
			Id:        23,
			Title:     "Article title",
			Subtitle:  "Article subtitle",
			Slug:      "article-slug",
			Body:      "Lorem ipsum dolar sit amet",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{Id: 42, Name: "tag one", Slug: "tag-one"},
				{Id: 69, Name: "tag two", Slug: "tag-two"},
			},
		},
	}, got, "should return article")
}

func testArticleRepoGetAllArticlesMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	tagQuery := `select id_, name_, slug_ from tags_`

	mockTagRows := mock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow(42, "tag one", "tag-one").
		AddRow(69, "tag two", "tag-two")

	mock.ExpectQuery(regexp.QuoteMeta(tagQuery)).WillReturnRows(mockTagRows)

	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour * 24)

	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ group by a.id_`

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
		"tag_ids_",
	}).AddRow(
		23,
		"Article title one",
		"Article subtitle one",
		"article-slug-one",
		"Lorem ipsum dolar sit amet one",
		createdAt,
		updatedAt,
		"42,69",
	).AddRow(
		42,
		"Article title two",
		"Article subtitle two",
		"article-slug-two",
		"Lorem ipsum dolar sit amet two",
		createdAt,
		updatedAt,
		"42,69",
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WillReturnRows(mockArticleRow)

	got, err := repo.GetAll()

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, err, "should not return error")

	require.Equal(t, &[]Article{
		{
			Id:        23,
			Title:     "Article title one",
			Subtitle:  "Article subtitle one",
			Slug:      "article-slug-one",
			Body:      "Lorem ipsum dolar sit amet one",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{Id: 42, Name: "tag one", Slug: "tag-one"},
				{Id: 69, Name: "tag two", Slug: "tag-two"},
			},
		},
		{
			Id:        42,
			Title:     "Article title two",
			Subtitle:  "Article subtitle two",
			Slug:      "article-slug-two",
			Body:      "Lorem ipsum dolar sit amet two",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{Id: 42, Name: "tag one", Slug: "tag-one"},
				{Id: 69, Name: "tag two", Slug: "tag-two"},
			},
		},
	}, got, "should return article")
}

func testArticleRepoGetArticleByAttrTagsScanError(t *testing.T, mock pgxmock.PgxPoolIface, repo ArticleRepository) {
	articleQuery := `select a.id_, a.title_, a.subtitle_, a.slug_, a.body_, a.created_at_, a.updated_at_, array_to_string(array_agg(distinct t.tag_id_), ',', '*') from articles_ a join article_tags_ t on a.id_ = t.article_id_ where a.slug_ = $1 group by a.id_`

	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Hour * 24)

	mockArticleRow := mock.NewRows([]string{
		"id_",
		"title_",
		"subtitle_",
		"slug_",
		"body_",
		"created_at_",
		"updated_at_",
		"tag_ids_",
	}).AddRow(
		23,
		"Article title",
		"Article subtitle",
		"article-slug",
		"Lorem ipsum dolar sit amet",
		createdAt,
		updatedAt,
		"42,69",
	)

	mock.
		ExpectQuery(regexp.QuoteMeta(articleQuery)).
		WithArgs("tag-slug").
		WillReturnRows(mockArticleRow)

	tagQuery := `select id_, name_, slug_ from tags_ where id_ = 42 or id_ = 69`

	mockBadTagRows := pgxmock.
		NewRows([]string{"id_", "name_", "slug_"}).
		AddRow("23", false, nil)

	mock.
		ExpectQuery(regexp.QuoteMeta(tagQuery)).
		WillReturnRows(mockBadTagRows)

	got, err := repo.GetByAttribute("slug", "tag-slug")

	mock.Reset()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, got, "should not return tags")

	require.Error(t, err, "should return error")
}
