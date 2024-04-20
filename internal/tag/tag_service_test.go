package tag

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTagData struct {
	mock.Mock
}

func (m *MockTagData) create(tag *Tag) (*Tag, error) {
	args := m.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagData) deleteById(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

var mockData = new(MockTagData)

func TestTagServiceDeleteById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"delete tag without error": testDeleteTagWithoutError,
		"delete tag with error":    testDeleteTagWithError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData)

			fn(t, service)
		})
	}
}

func testDeleteTagWithoutError(t *testing.T, service TagService) {
	mockData.On("deleteById", 23).Return(nil)

	err := service.DeleteById(23)

	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")
}

func testDeleteTagWithError(t *testing.T, service TagService) {
	mockData.On("deleteById", 42).Return(errors.New("data_error"))

	err := service.DeleteById(42)

	mockData.AssertExpectations(t)

	require.EqualError(t, err, "data_error", "should return error from data layer")
}
