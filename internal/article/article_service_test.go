package article

import (
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockData = new(MockArticleRepository)

var validate, _ = validation.NewValidator()

func TestArticleService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service ArticleService){
		"create article (success)":                      testArticleServiceCreateArticle,
		"create article (error - fails validation)":     testArticleServiceCreateFailsValidation,
		"create article (error - fail with no tags)":    testArticleServiceCreateArticleNoTags,
		"create article (error - repo error)":           testArticleServiceCreateArticleRepoError,
		"get all articles (success - multiple results)": testArticleServiceGetAllArticles,
		"get all articles (error)":                      testArticleServiceGetAllArticlesError,
		"get by slug (success)":                         testArticleServiceGetArticleBySlug,
		"get by slug (error)":                           testArticleServiceGetArticleBySlugError,
		"get many by attriute (success)":                testArticleServiceGetManyArticlesByTagSlug,
		"get many by attriute (error)":                  testArticleServiceGetManyArticlesByTagSlugError,
		"update article (success)":                      testArticleServiceUpdateArticle,
		"update article (error - fails validation)":     testArticleServiceUpdateFailsValidation,
		"update article (error)":                        testArticleServiceUpdateArticleError,
		"delete article by id (success)":                testArticleServiceDeleteArticleById,
		"delete article by id (error)":                  testArticleServiceDeleteArticleByIdError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewArticleService(mockData, validate)
			fn(t, service)
		})
	}
}

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

func testArticleServiceCreateArticle(t *testing.T, service ArticleService) {
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

	newArticle := ArticleNewRequestDto{
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Nil(t, err, "should not error")

	require.Equal(
		t,
		&mockCreatedArticle,
		createdArticle,
		"should return created article",
	)
}

func testArticleServiceCreateArticleNoTags(t *testing.T, service ArticleService) {
	articleWithNoTags := ArticleNewRequestDto{
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

func testArticleServiceCreateArticleRepoError(t *testing.T, service ArticleService) {
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

	newArticle := ArticleNewRequestDto{
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Nil(t, article, "should not return article")
	require.EqualError(t, err, "repo_error", "should return error")
}

func testArticleServiceDeleteArticleByIdError(t *testing.T, service ArticleService) {
	mockCall := mockData.On("DeleteById", 23).Return(errors.New("repo_error"))

	err := service.DeleteById(23)

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.EqualError(t, err, "repo_error", "should bubble up error from repo")
}

func testArticleServiceDeleteArticleById(t *testing.T, service ArticleService) {
	mockCall := mockData.On("DeleteById", 23).Return(nil)

	err := service.DeleteById(23)

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Nil(t, err, "should not return error")
}

func testArticleServiceGetAllArticles(t *testing.T, service ArticleService) {
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Nil(t, err, "should not return error")

	require.Equal(t, articles, &allArticles, "should return all articles")
}

func testArticleServiceGetAllArticlesError(t *testing.T, service ArticleService) {
	mockCall := mockData.On("GetAll").Return(&[]Article{}, errors.New("repo_error"))

	articles, err := service.GetAll()

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.EqualError(t, err, "repo_error", "should return error from repo")
	require.Empty(t, articles, "should return empty articles")
}

func testArticleServiceGetArticleBySlug(t *testing.T, service ArticleService) {
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

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

func testArticleServiceGetArticleBySlugError(t *testing.T, service ArticleService) {
	mockCall := mockData.
		On("GetByAttribute", "slug", "article-slug").
		Return(&Article{}, errors.New("repo_error"))

	gotArticle, err := service.GetByAttribute("slug", "article-slug")

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.EqualError(t, err, "repo_error", "should return error")
	require.Nil(t, gotArticle, "should not return an article")
}

func testArticleServiceUpdateArticle(t *testing.T, service ArticleService) {
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

	articleUpdate := ArticleUpdateRequestDto{
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

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

func testArticleServiceUpdateArticleError(t *testing.T, service ArticleService) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(42)

	articleUpdate := ArticleUpdateRequestDto{
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
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.EqualError(t, err, "repo_error", "should return error")
	require.Empty(t, updated, "should not return non-updated article")

}

func testArticleServiceCreateFailsValidation(t *testing.T, service ArticleService) {
	gotMissingFields, err := service.Create(&ArticleNewRequestDto{})

	require.Nil(t, gotMissingFields, "should not create article")

	missingFieldErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		missingFieldErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMissingFields, "should not return a tag")

	require.Equal(t, "required", missingFieldErrs["Title"], "should error for no Title")
	require.Equal(t, "required", missingFieldErrs["Subtitle"], "should error for no Subtitle")
	require.Equal(t, "required", missingFieldErrs["Slug"], "should error for no Slug")
	require.Equal(t, "required", missingFieldErrs["Body"], "should error for no Body")
	require.Equal(t, "required", missingFieldErrs["CreatedAt"], "should error for no CreatedAt")
	require.Equal(t, "required", missingFieldErrs["UpdatedAt"], "should error for no UpdatedAt")
	require.Equal(t, "required", missingFieldErrs["TagIds"], "should error for no TagIds")

	gotMaxValidations, err := service.Create(&ArticleNewRequestDto{
		Title:     "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345",
		Subtitle:  "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345",
		Slug:      "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    []int{1, 2, 3},
	})

	maxValidationErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		maxValidationErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMaxValidations, "should not return a tag")

	require.Equal(t, "max", maxValidationErrs["Title"], "should not allow long Title")
	require.Equal(t, "max", maxValidationErrs["Subtitle"], "should not allow long Subtitle")
	require.Equal(t, "max", maxValidationErrs["Slug"], "should not allow long Slug")

	gotMinValidations, err := service.Create(&ArticleNewRequestDto{
		Title:     "Some title",
		Subtitle:  "Some subtitle",
		Slug:      "a",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    []int{1, 2, 3},
	})

	minValidationErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		minValidationErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMinValidations, "should not return a tag")

	require.Equal(t, "min", minValidationErrs["Slug"], "should not allow short Slug")
}

func testArticleServiceUpdateFailsValidation(t *testing.T, service ArticleService) {
	gotMissingFields, err := service.Update(&ArticleUpdateRequestDto{})

	require.Nil(t, gotMissingFields, "should not update article")

	missingFieldErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		missingFieldErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMissingFields, "should not return a tag")

	require.Equal(t, "required", missingFieldErrs["Title"], "should error for no Title")
	require.Equal(t, "required", missingFieldErrs["Subtitle"], "should error for no Subtitle")
	require.Equal(t, "required", missingFieldErrs["Slug"], "should error for no Slug")
	require.Equal(t, "required", missingFieldErrs["Body"], "should error for no Body")
	require.Equal(t, "required", missingFieldErrs["CreatedAt"], "should error for no CreatedAt")
	require.Equal(t, "required", missingFieldErrs["UpdatedAt"], "should error for no UpdatedAt")
	require.Equal(t, "required", missingFieldErrs["TagIds"], "should error for no TagIds")

	gotMaxValidations, err := service.Update(&ArticleUpdateRequestDto{
		Title:     "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345",
		Subtitle:  "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345",
		Slug:      "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    []int{1, 2, 3},
	})

	maxValidationErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		maxValidationErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMaxValidations, "should not return a tag")

	require.Equal(t, "max", maxValidationErrs["Title"], "should not allow long Title")
	require.Equal(t, "max", maxValidationErrs["Subtitle"], "should not allow long Subtitle")
	require.Equal(t, "max", maxValidationErrs["Slug"], "should not allow long Slug")

	gotMinValidations, err := service.Update(&ArticleUpdateRequestDto{
		Title:     "Some title",
		Subtitle:  "Some subtitle",
		Slug:      "a",
		Body:      "Lorem ipsum dolar sit amet...",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    []int{1, 2, 3},
	})

	minValidationErrs := make(map[string]string)

	for _, v := range err.(validator.ValidationErrors) {
		minValidationErrs[v.Field()] = v.Tag()
	}

	require.Nil(t, gotMinValidations, "should not return a tag")

	require.Equal(t, "min", minValidationErrs["Slug"], "should not allow short Slug")
}

func testArticleServiceGetManyArticlesByTagSlug(t *testing.T, service ArticleService) {
	createdAt := time.Now()
	updatedAt := time.Now().Add(53)

	mockRepoArticles := []Article{
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

	mockCall := mockData.
		On("GetManyByAttribute", "tagSlug", "tag-one").
		Return(&mockRepoArticles, nil)

	gotArticle, err := service.GetManyByAttribute("tagSlug", "tag-one")

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]ArticleResponseDto{
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
	}, gotArticle, "should return article by slug")
}

func testArticleServiceGetManyArticlesByTagSlugError(t *testing.T, service ArticleService) {
	mockCall := mockData.
		On("GetManyByAttribute", "tagSlug", "tag-one").
		Return(&[]Article{}, errors.New("repo_error"))

	got, err := service.GetManyByAttribute("tagSlug", "tag-one")

	mockCall.Unset()
	if res := mockData.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	require.Empty(t, got, "should not return article(s)")
	require.EqualError(t, err, "repo_error", "should return repo error")
}
