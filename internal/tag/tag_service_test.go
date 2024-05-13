package tag

import (
	"errors"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockData = new(MockTagRepository)

func TestTagService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service TagServiceImpl){
		"update tag (success)":                              testTagServiceUpdateTag,
		"update tag (success - converts slug to lowercase)": testTagServiceUpdateConvertSlugToLowercase,
		"update tag (handle error from repo)":               testTagServiceUpdateTagRepoError,
		"get by slug (success - tag exists)":                testTagServiceGetBySlugTagExists,
		"get by slug (error - tag does not exist)":          testTagServiceGetBySlugTagDoesNotExist,
		"get all (multiple results)":                        testTagServiceGetAllTagsMultipleResults,
		"get all (no results)":                              testTagServiceGetAllTagsNoResults,
		"get all (single result)":                           testTagServiceGetAllTagsSingleResult,
		"get all (handle error from repo)":                  testTagServiceGetAllTagsRepoError,
		"create (success)":                                  testTagServiceCreateValidTag,
		"create (success - converts slug to lowercase)":     testTagServiceCreateConvertSlugToLowercase,
		"create (fail to create invalid tag)":               testTagServiceCreateInvalidTag,
		"create (fail to update invalid tag)":               testTagServiceUpdateInvalidTag,
		"create (fail to create existing tag)":              testTagServiceCreateExistingTag,
		"delete (success - delete tag by id)":               testTagServiceDeleteTagWithoutError,
		"delete (fail to delete tag by non-existent id)":    testTagServiceDeleteTagWithError,
	}

	var validate, err = pkg.NewValidator()
	if err != nil {
		t.Fatal("could not create validator", err.Error())
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			service := NewTagService(mockData, validate)

			fn(t, service)
		})
	}
}

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Create(tag *Tag) (*Tag, error) {
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

func (m *MockTagRepository) Exists(tag *Tag) (bool, error) {
	args := m.Called(tag)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTagRepository) GetAll() (*[]Tag, error) {
	args := m.Called()

	return args.Get(0).(*[]Tag), args.Error(1)
}

func (m *MockTagRepository) GetByAttribute(attr, slug string) (*Tag, error) {
	args := m.Called(attr, slug)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockTagRepository) Update(tag *Tag) (*Tag, error) {
	args := m.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
}

func testTagServiceUpdateConvertSlugToLowercase(t *testing.T, service TagServiceImpl) {
	mockRepoUpdate := mockData.On("Update", &Tag{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&Tag{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	got, err := service.Update(&TagUpdateRequestDto{
		Id:   42,
		Name: "tag name",
		Slug: "TaG-slUg",
	})

	mockRepoUpdate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")

	require.Equal(t, &TagResponseDto{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}, got, "should return updated tag")
}

func testTagServiceUpdateTag(t *testing.T, service TagServiceImpl) {
	mockRepoUpdate := mockData.On("Update", &Tag{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&Tag{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	got, err := service.Update(&TagUpdateRequestDto{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	})

	mockRepoUpdate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error out")

	require.Equal(t, &TagResponseDto{
		Id:   42,
		Name: "tag name",
		Slug: "tag-slug",
	}, got, "should return updated tag")
}

func testTagServiceDeleteTagWithoutError(t *testing.T, service TagServiceImpl) {
	mockRepoDeleteById := mockData.On("DeleteById", 23).Return(nil)

	got := service.DeleteById(23)

	mockRepoDeleteById.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, got, "should not error out")
}

func testTagServiceDeleteTagWithError(t *testing.T, service TagServiceImpl) {
	mockRepoDeleteById := mockData.On("DeleteById", 42).Return(errors.New("data_error"))

	got := service.DeleteById(42)

	mockRepoDeleteById.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, got, "data_error", "should return error from data layer")
}

func testTagServiceCreateConvertSlugToLowercase(t *testing.T, service TagServiceImpl) {
	mockRepoCreate := mockData.On("Create", &Tag{
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	got, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "tAg-SluG",
	})

	mockRepoCreate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error")

	require.Equal(t, &TagResponseDto{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, got, "should return tag with id")
}

func testTagServiceCreateValidTag(t *testing.T, service TagServiceImpl) {
	mockRepoCreate := mockData.On("Create", &Tag{
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	got, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "tag-slug",
	})

	mockRepoCreate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not error")

	require.Equal(t, &TagResponseDto{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, got, "should return tag with id")
}

func testTagServiceCreateInvalidTag(t *testing.T, service TagServiceImpl) {
	gotLongTagName, err := service.Create(&TagNewRequestDto{
		Name: "tag name that is longer than 50 characters so exceeds limit",
		Slug: "tag-slug",
	})

	require.NotNil(t, err, "should return error")
	require.Nil(t, gotLongTagName, "should not create tag")

	gotShortTagName, err := service.Create(&TagNewRequestDto{
		Name: "s",
		Slug: "tag-slug",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotShortTagName, "should not create tag")

	gotLongTagSlug, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "tag-slug-that-is-longer-than-50-characters-so-is-invalid",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotLongTagSlug, "should not create tag")

	gotShortTagSlug, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "s",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotShortTagSlug, "should not create tag")

	gotInvalidTagSlugWithSpecials, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "s%l&u*g",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotInvalidTagSlugWithSpecials, "should not create tag")

	gotInvalidTagSlugWithSpaces, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "s l u g",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotInvalidTagSlugWithSpaces, "should not create tag")
}

func testTagServiceUpdateInvalidTag(t *testing.T, service TagServiceImpl) {
	gotLongTagName, err := service.Update(&TagUpdateRequestDto{
		Id:   69,
		Name: "tag name that is longer than 50 characters so exceeds limit",
		Slug: "tag-slug",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotLongTagName, "should not update tag")

	gotShortTagName, err := service.Update(&TagUpdateRequestDto{
		Id:   69,
		Name: "s",
		Slug: "tag-slug",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotShortTagName, "should not update tag")

	gotLongTagSlug, err := service.Update(&TagUpdateRequestDto{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug-that-is-longer-than-50-characters-so-is-invalid",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotLongTagSlug, "should not update tag")

	gotShortTagSlug, err := service.Update(&TagUpdateRequestDto{
		Id:   69,
		Name: "tag name",
		Slug: "s",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotShortTagSlug, "should not update tag")

	gotInvalidTagSlugWithSpecials, err := service.Update(&TagUpdateRequestDto{
		Id:   1,
		Name: "tag name",
		Slug: "s%l&u*g",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotInvalidTagSlugWithSpecials, "should not update tag")

	gotInvalidTagSlugWithSpaces, err := service.Update(&TagUpdateRequestDto{
		Name: "tag name",
		Slug: "s l u g",
	})
	require.NotNil(t, err, "should return error")
	require.Nil(t, gotInvalidTagSlugWithSpaces, "should not create tag")
}

func testTagServiceCreateExistingTag(t *testing.T, service TagServiceImpl) {
	mockRepoCreate := mockData.On("Create", &Tag{
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(nil, errors.New("exists"))

	gotTag, err := service.Create(&TagNewRequestDto{
		Name: "tag name",
		Slug: "tag-slug",
	})

	mockRepoCreate.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, gotTag, "should not return tag")
	require.NotNil(t, err, "should return error")
}

func testTagServiceGetAllTagsNoResults(t *testing.T, service TagServiceImpl) {
	mockRepoGetAll := mockData.On("GetAll").Return(&[]Tag{}, nil)

	got, err := service.GetAll()

	mockRepoGetAll.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Empty(t, got, "should not return any tags")
}

func testTagServiceGetAllTagsMultipleResults(t *testing.T, service TagServiceImpl) {
	mockRepoGetAll := mockData.On("GetAll").Return(&[]Tag{
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

	got, err := service.GetAll()

	mockRepoGetAll.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]TagResponseDto{
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
	}, got, "should return all tags")
}

func testTagServiceGetAllTagsSingleResult(t *testing.T, service TagServiceImpl) {
	mockRepoGetAll := mockData.On("GetAll").Return(&[]Tag{
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, nil)

	tags, err := service.GetAll()

	mockRepoGetAll.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &[]TagResponseDto{
		{
			Id:   69,
			Name: "tagname3",
			Slug: "tag-slug-3",
		},
	}, tags, "should return all tags")
}

func testTagServiceGetBySlugTagExists(t *testing.T, service TagServiceImpl) {
	mockRepoGetByAttribute := mockData.On("GetByAttribute", "slug", "tag-slug").Return(&Tag{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	got, err := service.GetByAttribute("slug", "tag-slug")

	mockRepoGetByAttribute.Unset()
	mockData.AssertExpectations(t)

	require.Nil(t, err, "should not return error")
	require.Equal(t, &TagResponseDto{
		Id:   69,
		Name: "tag name",
		Slug: "tag-slug",
	}, got, "should return tag")
}

func testTagServiceGetBySlugTagDoesNotExist(t *testing.T, service TagServiceImpl) {
	mockRepoGetByAttribute := mockData.On("GetByAttribute", "slug", "tag-slug").Return(nil, errors.New("data_error"))

	got, err := service.GetByAttribute("slug", "tag-slug")

	mockRepoGetByAttribute.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "data_error", "should return error from data layer")
	require.Nil(t, got, "should not return tag")
}

func testTagServiceGetAllTagsRepoError(t *testing.T, service TagServiceImpl) {
	mockRepoGetAll := mockData.On("GetAll").Return(&[]Tag{}, errors.New("getall_repo_error"))

	got, err := service.GetAll()

	mockRepoGetAll.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "getall_repo_error", "should return error from repo method call")
	require.Nil(t, got, "should not return any tags")
}

func testTagServiceUpdateTagRepoError(t *testing.T, service TagServiceImpl) {
	mockCall := mockData.On("Update", &Tag{
		Id:   23,
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&Tag{}, errors.New("update_repo_error"))

	got, err := service.Update(&TagUpdateRequestDto{
		Id:   23,
		Name: "tag name",
		Slug: "tag-slug",
	})

	mockCall.Unset()
	mockData.AssertExpectations(t)

	require.EqualError(t, err, "update_repo_error", "should return error")
	require.Nil(t, got, "should not return/update a tag")
}
