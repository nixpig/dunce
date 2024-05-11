package article

import (
	"errors"
	"testing"
	"time"

	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(article *ArticleNew) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleRepository) GetAll() (*[]Article, error) {
	args := m.Called()

	return args.Get(0).(*[]Article), args.Error(1)
}

func (m *MockArticleRepository) GetByAttribute(attr, value string) (*Article, error) {
	args := m.Called(attr, value)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleRepository) Update(article *UpdateArticle) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleRepository) DeleteById(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockArticleRepository) Exists(article *Article) (bool, error) {
	args := m.Called(article)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockArticleRepository) GetManyByAttribute(attr, val string) (*[]Article, error) {
	args := m.Called(attr, val)

	return args.Get(0).(*[]Article), args.Error(1)

}

var mockData = new(MockArticleRepository)

var validate, _ = pkg.NewValidator()

func TestArticleServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service ArticleServiceImpl){
		// create
		"create article":                       testServiceCreateArticle,
		"fail to create article with no tags":  testServiceCreateArticleNoTags,
		"fail to create article on repo error": testServiceCreateArticleRepoError,

		// read
		"get all articles":       testServiceGetAllArticles,
		"get all articles error": testServiceGetAllArticlesError,
		"get by slug":            testServiceGetArticleBySlug,
		"get by slug error":      testServiceGetArticleBySlugError,

		// update
		"update article success": testServiceUpdateArticle,
		"update article error":   testServiceUpdateArticleError,

		// delete
		"delete article by id":       testServiceDeleteArticleById,
		"delete article by id error": testServiceDeleteArticleByIdError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewArticleService(mockData, validate)
			fn(t, service)
		})
	}
}

func testServiceCreateArticle(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now()

	mockArticleCall := ArticleNew{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds: []int{
			1,
			2,
		},
	}

	newArticle := ArticleRequestDto{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds: []int{
			1,
			2,
		},
	}

	mockCreatedArticle := ArticleResponseDto{
		Id:        42,
		Title:     newArticle.Title,
		Subtitle:  newArticle.Subtitle,
		Slug:      newArticle.Slug,
		Body:      newArticle.Body,
		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
		Tags: []tag.Tag{
			{
				Id:   1,
				Name: "tag one",
				Slug: "slug-one",
			},
			{
				Id:   2,
				Name: "tag two",
				Slug: "slug-two",
			},
		},
	}

	mockRepoArticleResponse := Article{
		Id:        42,
		Title:     newArticle.Title,
		Subtitle:  newArticle.Subtitle,
		Slug:      newArticle.Slug,
		Body:      newArticle.Body,
		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
		Tags: []tag.Tag{
			{
				Id:   1,
				Name: "tag one",
				Slug: "slug-one",
			},
			{
				Id:   2,
				Name: "tag two",
				Slug: "slug-two",
			},
		},
	}

	mockCallCreate := mockData.On("Create", &mockArticleCall).Return(&mockRepoArticleResponse, nil)

	createdArticle, err := service.Create(&newArticle)

	mockCallCreate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error")

	require.Equal(
		t,
		&mockCreatedArticle,
		createdArticle,
		"should return created article",
	)
}

func testServiceCreateArticleNoTags(t *testing.T, service ArticleServiceImpl) {
	articleWithNoTags := ArticleRequestDto{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(23),
		TagIds:    []int{},
	}

	article, err := service.Create(&articleWithNoTags)

	require.Nil(t, article, "no article should be returned")
	require.EqualError(t, err, "article must have at least one tag", "should return error indicating article requires one or more tags")
}

func testServiceCreateArticleRepoError(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now()

	mockArticleData := ArticleNew{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{1, 2},
	}

	mockCall := mockData.On("Create", &mockArticleData).Return(&Article{}, errors.New("repo_error"))

	newArticle := ArticleRequestDto{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{1, 2},
	}
	article, err := service.Create(&newArticle)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, article, "should not return article")
	require.EqualError(t, err, "repo_error", "should return error")
}

func testServiceDeleteArticleByIdError(t *testing.T, service ArticleServiceImpl) {
	mockCall := mockData.On("DeleteById", 23).Return(errors.New("repo_error"))

	err := service.DeleteById(23)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "repo_error", "should bubble up error from repo")
}

func testServiceDeleteArticleById(t *testing.T, service ArticleServiceImpl) {
	mockCall := mockData.On("DeleteById", 23).Return(nil)

	err := service.DeleteById(23)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
}

func testServiceGetAllArticles(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(42)

	allArticles := []ArticleResponseDto{
		{
			Title:     "article one title",
			Subtitle:  "article one subtitle",
			Slug:      "article-one-slug",
			Body:      "article one body content",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{
					Id:   23,
					Name: "tag one",
					Slug: "tag-one",
				},
			},
		},
		{
			Title:     "article two title",
			Subtitle:  "article two subtitle",
			Slug:      "article-two-slug",
			Body:      "article two body content",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{
					Id:   23,
					Name: "tag one",
					Slug: "tag-one",
				},
			},
		},
	}

	mockAllArticles := []Article{
		{
			Title:     "article one title",
			Subtitle:  "article one subtitle",
			Slug:      "article-one-slug",
			Body:      "article one body content",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{
					Id:   23,
					Name: "tag one",
					Slug: "tag-one",
				},
			},
		},
		{
			Title:     "article two title",
			Subtitle:  "article two subtitle",
			Slug:      "article-two-slug",
			Body:      "article two body content",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Tags: []tag.Tag{
				{
					Id:   23,
					Name: "tag one",
					Slug: "tag-one",
				},
			},
		},
	}

	mockCall := mockData.On("GetAll").Return(&mockAllArticles, nil)

	articles, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")

	require.Equal(t, articles, &allArticles, "should return all articles")
}

func testServiceGetAllArticlesError(t *testing.T, service ArticleServiceImpl) {
	mockCall := mockData.On("GetAll").Return(&[]Article{}, errors.New("repo_error"))

	articles, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "repo_error", "should return error from repo")
	require.Empty(t, articles, "should return empty articles")
}

func testServiceGetArticleBySlug(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(53)

	mockRepoArticle := Article{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags: []tag.Tag{
			{
				Id:   23,
				Name: "tag one",
				Slug: "tag-one",
			},
		},
	}

	mockCall := mockData.
		On("GetByAttribute", "slug", "article-slug").
		Return(&mockRepoArticle, nil)

	gotArticle, err := service.GetByAttribute("slug", "article-slug")

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &ArticleResponseDto{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags: []tag.Tag{
			{
				Id:   23,
				Name: "tag one",
				Slug: "tag-one",
			},
		},
	}, gotArticle, "should return article by slug")
}

func testServiceGetArticleBySlugError(t *testing.T, service ArticleServiceImpl) {
	mockCall := mockData.
		On("GetByAttribute", "slug", "article-slug").
		Return(&Article{}, errors.New("repo_error"))

	gotArticle, err := service.GetByAttribute("slug", "article-slug")

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "repo_error", "should return error")
	require.Nil(t, gotArticle, "should not return an article")
}

func testServiceUpdateArticle(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now()

	mockUpdateArticle := UpdateArticle{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{23},
	}

	articleUpdate := UpdateArticleRequestDto{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{23},
	}

	mockUpdateArticleRepo := Article{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags: []tag.Tag{
			{
				Id:   23,
				Name: "foo",
				Slug: "bar-baz",
			},
		},
	}

	mockCall := mockData.
		On("Update", &mockUpdateArticle).
		Return(&mockUpdateArticleRepo, nil)

	updated, err := service.Update(&articleUpdate)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &ArticleResponseDto{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Tags: []tag.Tag{
			{
				Id:   23,
				Name: "foo",
				Slug: "bar-baz",
			},
		},
	}, updated, "should return updated article")
}

func testServiceUpdateArticleError(t *testing.T, service ArticleServiceImpl) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(42)

	articleUpdate := UpdateArticleRequestDto{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{23},
	}

	mockUpdateArticle := UpdateArticle{
		Title:     "article one title",
		Subtitle:  "article one subtitle",
		Slug:      "article-one-slug",
		Body:      "article one body content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TagIds:    []int{23},
	}

	mockCall := mockData.
		On("Update", &mockUpdateArticle).
		Return(&Article{}, errors.New("repo_error"))

	updated, err := service.Update(&articleUpdate)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "repo_error", "should return error")
	require.Empty(t, updated, "should not return non-updated article")

}
