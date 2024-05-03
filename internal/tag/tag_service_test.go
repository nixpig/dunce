package tag

import (
	"errors"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Create(tag *TagData) (*Tag, error) {
	args := m.Called(tag)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagRepository) DeleteById(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockTagRepository) Exists(tag *TagData) (bool, error) {
	args := m.Called(tag)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTagRepository) GetAll() (*[]Tag, error) {
	args := m.Called()

	return args.Get(0).(*[]Tag), args.Error(1)
}

func (m *MockTagRepository) GetBySlug(slug string) (*Tag, error) {
	args := m.Called(slug)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagRepository) Update(tag *Tag) (*Tag, error) {
	args := m.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
}

var mockData = new(MockTagRepository)

var validate, _ = pkg.NewValidator()

func TestTagService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagService){
		"update tag":                            testServiceUpdateTag,
		"update tag (error from repo)":          testServiceUpdateTagRepoError,
		"get by slug (tag exists)":              testServiceGetBySlugTagExists,
		"get by slug (tag does not exist)":      testServiceGetBySlugTagDoesNotExist,
		"get all (multiple results)":            testServiceGetAllTagsMultipleResults,
		"get all (no results)":                  testServiceGetAllTagsNoResults,
		"get all (single result)":               testServiceGetAllTagsSingleResult,
		"get all (error from repo)":             testServiceGetAllTagsRepoError,
		"successfully create valid tag":         testTagServiceCreateValidTag,
		"fail to create invalid tag":            testTagServiceCreateInvalidTag,
		"fail to update invalid tag":            testTagServiceUpdateInvalidTag,
		"fail to create existing tag":           testTagServiceCreateExistingTag,
		"successfully delete tag by id":         testTagServiceDeleteTagWithoutError,
		"fail to delete tag by non-existent id": testTagServiceDeleteTagWithError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate, pkg.NewLogger())

			fn(t, service)
		})
	}
}

func testServiceUpdateTag(t *testing.T, service TagService) {
	tag := Tag{
		Id: 42,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug",
		},
	}

	mockCallUpdate := mockData.On("Update", &tag).Return(&tag, nil)

	updatedTag, err := service.Update(&tag)

	mockCallUpdate.Unset()

	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")
	require.Equal(t, &tag, updatedTag, "should return updated tag")
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
	newTag := TagData{Name: "tag name", Slug: "tag-slug"}

	mockCreatedTag := Tag{
		Id: 69,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug",
		},
	}

	mockCallCreate := mockData.On("Create", &newTag).Return(&mockCreatedTag, nil)

	createdTag, err := service.Create(&newTag)

	mockCallCreate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error")

	require.Equal(t, &Tag{
		Id: 69,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug",
		},
	}, createdTag, "should return tag with id")
}

func testTagServiceCreateInvalidTag(t *testing.T, service TagService) {
	longTagName := TagData{
		Name: "tag name that is longer than 50 characters so exceeds limit",
		Slug: "tag-slug",
	}

	createdLongTagName, err := service.Create(&longTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdLongTagName, "should not create tag")

	shortTagName := TagData{
		Name: "s",
		Slug: "tag-slug",
	}

	createdShortTagName, err := service.Create(&shortTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdShortTagName, "should not create tag")

	longTagSlug := TagData{
		Name: "tag name",
		Slug: "tag-slug-that-is-longer-than-50-characters-so-is-invalid",
	}

	createdLongTagSlug, err := service.Create(&longTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdLongTagSlug, "should not create tag")

	shortTagSlug := TagData{
		Name: "tag name",
		Slug: "s",
	}

	createdShortTagSlug, err := service.Create(&shortTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdShortTagSlug, "should not create tag")

	incorrectFormatSlugUppercase := TagData{
		Name: "tag name",
		Slug: "Sslug",
	}

	createdIncorrectTagSlug, err := service.Create(&incorrectFormatSlugUppercase)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdIncorrectTagSlug, "should not create tag")

	incorrectFormatSlugNonSlug := TagData{Name: "tag name", Slug: "s%l&u*g"}

	createdIncorrectTagSlugNonSlug, err := service.Create(&incorrectFormatSlugNonSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, createdIncorrectTagSlugNonSlug, "should not create tag")
}

func testTagServiceUpdateInvalidTag(t *testing.T, service TagService) {
	longTagName := Tag{Id: 69,
		TagData: TagData{
			Name: "tag name that is longer than 50 characters so exceeds limit",
			Slug: "tag-slug",
		},
	}

	UpdatedLongTagName, err := service.Update(&longTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, UpdatedLongTagName, "should not update tag")

	shortTagName := Tag{Id: 69,
		TagData: TagData{
			Name: "s",
			Slug: "tag-slug",
		},
	}

	updatedShortTagName, err := service.Update(&shortTagName)
	require.NotNil(t, err, "should return error")
	require.Nil(t, updatedShortTagName, "should not update tag")

	longTagSlug := Tag{Id: 69,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug-that-is-longer-than-50-characters-so-is-invalid",
		},
	}

	updatedLongTagSlug, err := service.Update(&longTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, updatedLongTagSlug, "should not update tag")

	shortTagSlug := Tag{Id: 69,
		TagData: TagData{
			Name: "tag name",
			Slug: "s",
		},
	}

	updatedShortTagSlug, err := service.Update(&shortTagSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, updatedShortTagSlug, "should not update tag")

	incorrectFormatSlugUppercase := Tag{Id: 1, TagData: TagData{Name: "tag name", Slug: "Sslug"}}

	updatedIncorrectTagSlug, err := service.Update(&incorrectFormatSlugUppercase)
	require.NotNil(t, err, "should return error")
	require.Nil(t, updatedIncorrectTagSlug, "should not update tag")

	incorrectFormatSlugNonSlug := Tag{Id: 1, TagData: TagData{Name: "tag name", Slug: "s%l&u*g"}}

	updatedIncorrecTagSlugNonSlug, err := service.Update(&incorrectFormatSlugNonSlug)
	require.NotNil(t, err, "should return error")
	require.Nil(t, updatedIncorrecTagSlugNonSlug, "should not update tag")
}

func testTagServiceCreateExistingTag(t *testing.T, service TagService) {
	newTag := TagData{Name: "tag name", Slug: "tag-slug"}

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
			Id: 23,
			TagData: TagData{
				Name: "tagname1",
				Slug: "tag-slug-1",
			},
		},
		{
			Id: 42,
			TagData: TagData{
				Name: "tagname2",
				Slug: "tag-slug-2",
			},
		},
		{
			Id: 69,
			TagData: TagData{
				Name: "tagname3",
				Slug: "tag-slug-3",
			},
		},
	}, nil)

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Tag{
		{
			Id: 23,
			TagData: TagData{
				Name: "tagname1",
				Slug: "tag-slug-1",
			},
		},
		{
			Id: 42,
			TagData: TagData{
				Name: "tagname2",
				Slug: "tag-slug-2",
			},
		},
		{
			Id: 69,
			TagData: TagData{
				Name: "tagname3",
				Slug: "tag-slug-3",
			},
		},
	}, tags, "should return all tags")
}

func testServiceGetAllTagsSingleResult(t *testing.T, service TagService) {
	mockCall := mockData.On("GetAll").Return(&[]Tag{
		{
			Id: 69,
			TagData: TagData{
				Name: "tagname3",
				Slug: "tag-slug-3",
			},
		},
	}, nil)

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]Tag{
		{
			Id: 69,
			TagData: TagData{
				Name: "tagname3",
				Slug: "tag-slug-3",
			},
		},
	}, tags, "should return all tags")
}

func testServiceGetBySlugTagExists(t *testing.T, service TagService) {
	mockCall := mockData.On("GetBySlug", "tag-slug").Return(&Tag{
		Id: 69,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug",
		},
	}, nil)

	tag, err := service.GetBySlug("tag-slug")

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &Tag{
		Id: 69,
		TagData: TagData{

			Name: "tag name",
			Slug: "tag-slug",
		},
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

func testServiceGetAllTagsRepoError(t *testing.T, service TagService) {
	mockCall := mockData.On("GetAll").Return(&[]Tag{}, errors.New("getall_repo_error"))

	tags, err := service.GetAll()

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "getall_repo_error", "should return error from repo method call")
	require.Nil(t, tags, "should not return any tags")
}

func testServiceUpdateTagRepoError(t *testing.T, service TagService) {
	mockTagUpdate := Tag{
		Id: 23,
		TagData: TagData{
			Name: "tag name",
			Slug: "tag-slug",
		},
	}

	mockCall := mockData.On("Update", &mockTagUpdate).Return(&Tag{}, errors.New("update_repo_error"))

	tag, err := service.Update(&mockTagUpdate)

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "update_repo_error", "should return error")
	require.Nil(t, tag, "should not return/update a tag")
}
