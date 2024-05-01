package article

import (
	"testing"
	"time"

	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockArticleData struct {
	mock.Mock
}

func (m *MockArticleData) Create(article *Article) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleData) GetAll() (*[]Article, error) {
	args := m.Called()

	return args.Get(0).(*[]Article), args.Error(1)
}

func (m *MockArticleData) GetBySlug(slug string) (*Article, error) {
	args := m.Called(slug)

	return args.Get(0).(*Article), args.Error(1)
}

func (m *MockArticleData) Update(article *Article) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

var mockData = new(MockArticleData)

var validate, _ = validation.NewValidator()

func TestArticleServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service ArticleService){
		"create article": testServiceCreateArticle,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewArticleService(mockData, validate, logging.NewLogger())
			fn(t, service)
		})
	}
}

func testServiceCreateArticle(t *testing.T, service ArticleService) {
	newArticle := NewArticle("article title", "article subtitle", "article-slug", "article body content", time.Now(), time.Now().Add(23), []int{})

	mockCreatedArticle := NewArticleWithId(
		42,
		newArticle.Title,
		newArticle.Subtitle,
		newArticle.Slug,
		newArticle.Body,
		newArticle.CreatedAt,
		newArticle.UpdatedAt,
		[]int{},
	)

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
