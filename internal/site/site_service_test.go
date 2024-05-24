package site

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockRepo = new(MockSiteRepository)

func TestSiteService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service SiteService){
		"test create site service kv": testSiteServiceCreateKv,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {

			service := NewSiteService(mockRepo)

			fn(t, service)
		})
	}
}

type MockSiteRepository struct {
	mock.Mock
}

func (s *MockSiteRepository) Create(key, value string) (*SiteKv, error) {
	args := s.Called(key, value)

	return args.Get(0).(*SiteKv), args.Error(1)
}

func testSiteServiceCreateKv(t *testing.T, service SiteService) {
	mockSiteRepositoryCreate := mockRepo.
		On("Create", "some name", "some description").
		Return(&SiteKv{Key: "some name", Value: "some description"}, nil)

	got, err := service.Create("some name", "some description")

	require.NoError(t, err, "should not return error")
	require.Equal(t, &SiteKv{
		Key:   "some name",
		Value: "some description",
	}, got, "should return k/v pair")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("should call through to repo")
	}

	mockSiteRepositoryCreate.Unset()
}
