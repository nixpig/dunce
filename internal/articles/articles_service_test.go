package articles

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockArticleData struct {
	mock.Mock
}

func (m *MockArticleData) create(article *Article) (*Article, error) {
	args := m.Called(article)

	return args.Get(0).(*Article), args.Error(1)
}

var mockData = new(MockArticleData)

func TestArticleServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service ArticleService){
		"create article": testServiceCreateArticle,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewArticleService(mockData)
			fn(t, service)
		})
	}
}

func testServiceCreateArticle(t *testing.T, service ArticleService) {
	newArticle := NewArticle("article title", "article subtitle", "article-slug", "article body content", time.Now(), time.Now().Add(23))

	mockCreatedArticle := NewArticleWithId(
		42,
		newArticle.Title,
		newArticle.Subtitle,
		newArticle.Slug,
		newArticle.Body,
		newArticle.CreatedAt,
		newArticle.UpdatedAt,
	)

	mockCallCreate := mockData.On("create", &newArticle).Return(&mockCreatedArticle, nil)

	createdArticle, err := service.create(&newArticle)

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
