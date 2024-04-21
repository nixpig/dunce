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

func (m *MockTagData) exists(tag *Tag) (bool, error) {
	args := m.Called(tag)

	return args.Get(0).(bool), args.Error(1)
}

var mockData = new(MockTagData)

func TestTagServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"successfully create valid tag": testTagServiceCreateValidTag,
		"fail to create invalid tag":    testTagServiceCreateInvalidTag,
		"fail to create existing tag":   testTagServiceCreateExistingTag,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData)

			fn(t, service)
		})
	}
}

func TestTagServiceDeleteById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"successfully delete tag by id":         testTagServiceDeleteTagWithoutError,
		"fail to delete tag by non-existent id": testTagServiceDeleteTagWithError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData)

			fn(t, service)
		})
	}
}

func testTagServiceDeleteTagWithoutError(t *testing.T, service TagService) {
	mockData.On("deleteById", 23).Return(nil)

	err := service.DeleteById(23)

	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")
}

func testTagServiceDeleteTagWithError(t *testing.T, service TagService) {
	mockData.On("deleteById", 42).Return(errors.New("data_error"))

	err := service.DeleteById(42)

	mockData.AssertExpectations(t)

	require.EqualError(t, err, "data_error", "should return error from data layer")
}

func testTagServiceCreateValidTag(t *testing.T, service TagService) {
	newTag := NewTag("tag name", "tag-slug")

	mockCreatedTag := NewTagWithId(69, "tag name", "tag-slug")

	mockData.On("create", &newTag).Return(&mockCreatedTag, nil)
	mockData.On("exists", &newTag).Return(false, nil)

	createdTag, err := service.Create(&newTag)

	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error")

	require.Equal(t, &Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, createdTag, "should return tag with id")
}

func testTagServiceCreateInvalidTag(t *testing.T, service TagService) {
	longTagName := NewTagWithId(
		69,
		"tag name that is longer than 50 characters so exceeds limit",
		"tag-slug",
	)

	createdLongTagName, err := service.Create(&longTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdLongTagName, "should not create tag")

	shortTagName := NewTagWithId(
		69,
		"sn",
		"tag-slug",
	)

	createdShortTagName, err := service.Create(&shortTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdShortTagName, "should not create tag")

	longTagSlug := NewTagWithId(
		69,
		"tag name",
		"tag-slug-that-is-longer-than-50-characters-so-is-invalid",
	)

	createdLongTagSlug, err := service.Create(&longTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdLongTagSlug, "should not create tag")

	shortTagSlug := NewTagWithId(
		69,
		"tag name",
		"st",
	)

	createdShortTagSlug, err := service.Create(&shortTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdShortTagSlug, "should not create tag")
}

func testTagServiceCreateExistingTag(t *testing.T, service TagService) {
	newTag := NewTagWithId(42, "tag name", "tag-slug")

	mockData.On("exists", &newTag).Return(true, nil)

	_, err := service.Create(&newTag)

	mockData.AssertExpectations(t)

	require.EqualError(t, err, "tag name and/or slug already exists", "should return error")
}
