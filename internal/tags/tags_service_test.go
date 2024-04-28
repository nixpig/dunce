package tags

import (
	"errors"
	"regexp"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTagData struct {
	mock.Mock
}

func (m *MockTagData) Create(tag *Tag) (*Tag, error) {
	args := m.Called(tag)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagData) DeleteById(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockTagData) Exists(tag *Tag) (bool, error) {
	args := m.Called(tag)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTagData) GetAll() (*[]Tag, error) {
	args := m.Called()

	return args.Get(0).(*[]Tag), args.Error(1)
}

func (m *MockTagData) GetBySlug(slug string) (*Tag, error) {
	args := m.Called(slug)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagData) Update(tag *Tag) (*Tag, error) {
	args := m.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
}

var mockData = new(MockTagData)

var validate, _ = validation.NewValidator()

func TestTagServiceUpdate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"update tag": testServiceUpdateTag,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate)

			fn(t, service)
		})
	}
}

func testServiceUpdateTag(t *testing.T, service TagService) {
	tag := NewTagWithId(42, "tag name", "tag-slug")

	mockCallUpdate := mockData.On("Update", &tag).Return(&tag, nil)

	updatedTag, err := service.Update(&tag)

	mockCallUpdate.Unset()

	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")
	require.Equal(t, &tag, updatedTag, "should return updated tag")
}

func TestTagServiceGetBySlug(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"get by slug (tag exists)":         testServiceGetBySlugTagExists,
		"get by slug (tag does not exist)": testServiceGetBySlugTagDoesNotExist,
	}
	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate)

			fn(t, service)
		})
	}
}

// update

func TestTagServiceGetAll(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"get all (multiple results)": testServiceGetAllTagsMultipleResults,
		"get all (no results)":       testServiceGetAllTagsNoResults,
		"get all (single result)":    testServiceGetAllTagsSingleResult,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate)

			fn(t, service)
		})
	}
}

func TestTagServiceCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"successfully create valid tag": testTagServiceCreateValidTag,
		"fail to create invalid tag":    testTagServiceCreateInvalidTag,
		"fail to create existing tag":   testTagServiceCreateExistingTag,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate)

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
			service := NewTagService(mockData, validate)

			fn(t, service)
		})
	}
}

func testTagServiceDeleteTagWithoutError(t *testing.T, service TagService) {
	mockCall := mockData.On("DeleteById", 23).Return(nil)

	err := service.DeleteById(23)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")
}

func testTagServiceDeleteTagWithError(t *testing.T, service TagService) {
	mockCall := mockData.On("DeleteById", 42).Return(errors.New("data_error"))

	err := service.DeleteById(42)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "data_error", "should return error from data layer")
}

func testTagServiceCreateValidTag(t *testing.T, service TagService) {
	newTag := NewTag("tag name", "tag-slug")

	mockCreatedTag := NewTagWithId(69, "tag name", "tag-slug")

	mockCallCreate := mockData.On("Create", &newTag).Return(&mockCreatedTag, nil)

	createdTag, err := service.Create(&newTag)

	mockCallCreate.Unset()
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
		"s",
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
		"s",
	)

	createdShortTagSlug, err := service.Create(&shortTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdShortTagSlug, "should not create tag")
}

func testTagServiceCreateExistingTag(t *testing.T, service TagService) {
	newTag := NewTagWithId(42, "tag name", "tag-slug")

	mockCall := mockData.On("Create", &newTag).Return(nil, errors.New("exists"))

	createdTag, err := service.Create(&newTag)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, createdTag, "should not return tag")
	require.NotNil(t, err, "should return error")
}

func testServiceGetAllTagsNoResults(t *testing.T, service TagService) {
	mockCall := mockData.On("GetAll").Return(&[]Tag{}, nil)

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Empty(t, tags, "should not return any tags")
}

func testServiceGetAllTagsMultipleResults(t *testing.T, service TagService) {
	mockCall := mockData.On("GetAll").Return(&[]Tag{
		{
			Id:   23,
			Name: "tagname1",
			Slug: "tag-slug-1",
		},
		{
			Id:   42,
			Name: "tagname2",
			Slug: "tag-slug-2",
		},
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, nil)

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Tag{
		{
			Id:   23,
			Name: "tagname1",
			Slug: "tag-slug-1",
		},
		{
			Id:   42,
			Name: "tagname2",
			Slug: "tag-slug-2",
		},
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, tags, "should return all tags")
}

func testServiceGetAllTagsSingleResult(t *testing.T, service TagService) {
	mockCall := mockData.On("GetAll").Return(&[]Tag{
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, nil)

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Tag{
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, tags, "should return all tags")
}

func testServiceGetBySlugTagExists(t *testing.T, service TagService) {
	mockCall := mockData.On("GetBySlug", "tag-slug").Return(&Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	tag, err := service.GetBySlug("tag-slug")

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, tag, "should return tag")
}

func testServiceGetBySlugTagDoesNotExist(t *testing.T, service TagService) {
	mockCall := mockData.On("GetBySlug", "tag-slug").Return(nil, errors.New("data_error"))

	tag, err := service.GetBySlug("tag-slug")

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "data_error", "should return error from data layer")
	require.Nil(t, tag, "should not return tag")
}

func ValidateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9\\-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
