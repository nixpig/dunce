package article

import (
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

func (m *MockArticleRepository) GetBySlug(slug string) (*Article, error) {
	args := m.Called(slug)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleRepository) Update(article *Article) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleRepository) DeleteById(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockArticleRepository) Exists(article *ArticleNew) (bool, error) {
	args := m.Called(article)

	return args.Get(0).(bool), args.Error(1)
}

var mockData = new(MockArticleRepository)

var validate, _ = pkg.NewValidator()

func TestArticleServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service ArticleService){
		"create article": testServiceCreateArticle,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewArticleService(mockData, validate, pkg.NewLogger())
			fn(t, service)
		})
	}
}

func testServiceCreateArticle(t *testing.T, service ArticleService) {
	newArticle := ArticleNew{
		Title:     "article title",
		Subtitle:  "article subtitle",
		Slug:      "article-slug",
		Body:      "article body content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(23),
		TagIds:    []int{},
	}

	mockCreatedArticle := Article{
		Id:        42,
		Title:     newArticle.Title,
		Subtitle:  newArticle.Subtitle,
		Slug:      newArticle.Slug,
		Body:      newArticle.Body,
		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
		Tags: []tag.Tag{
			{
				Id: 1,
				TagData: tag.TagData{
					Name: "tag one",
					Slug: "slug-one",
				},
			},
			{
				Id: 2,
				TagData: tag.TagData{
					Name: "tag two",
					Slug: "slug-two",
				},
			},
		},
	}

	mockCallCreate := mockData.On("Create", &newArticle).Return(&mockCreatedArticle, nil)

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
