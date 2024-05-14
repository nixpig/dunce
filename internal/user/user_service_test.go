package user

import (
	"errors"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockRepo = new(MockUserRepo)

func TestUserService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service UserService){
		"get all users (success - multiple results)": testUserServiceGetAllMultiple,
		"get all users (success - zero results)":     testUserServiceGetAllZero,
		"get all users (success - single result)":    testUserServiceGetAllSingle,
		"get all users (error - repo)":               testUserServiceRepoError,
		"get by attribute (success)":                 testUserServiceGetByAttr,
		"get by attribute (error - repo)":            testUserServiceGetByAttrRepoError,
		"user exists (success - true)":               testUserServiceUserExistsTrue,
		"user exists (success - false)":              testUserServiceUserExistsFalse,
		"user exists (error - false)":                testUserServiceUserExistsError,
		"delete user by id (success)":                testUserServiceDeleteById,
		"delete user by id (error)":                  testUserServiceDeleteByIdError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			validator, err := pkg.NewValidator()
			if err != nil {
				t.Fatal("unable to construct validator")
			}

			service := NewUserService(mockRepo, validator)

			fn(t, service)
		})
	}
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *User) (*User, error) {
	args := m.Called(user)

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) DeleteById(id uint) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *MockUserRepo) Exists(username string) (bool, error) {
	args := m.Called(username)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepo) GetAll() (*[]User, error) {
	args := m.Called()

	return args.Get(0).(*[]User), args.Error(1)
}

func (m *MockUserRepo) GetByAttribute(attr, value string) (*User, error) {
	args := m.Called(attr, value)

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) Update(user *User) (*User, error) {
	args := m.Called(user)

	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) GetPasswordByUsername(username string) (string, error) {
	args := m.Called(username)

	return args.Get(0).(string), args.Error(1)
}

func testUserServiceGetAllMultiple(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{
		{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		},
		{
			Id:       42,
			Username: "johndoe",
			Email:    "john@example.net",
		},
	}, nil)

	users, err := service.GetAll()

	require.NoError(t, err, "should not return an error")

	require.Equal(t, &[]UserResponseDto{
		{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		},
		{
			Id:       42,
			Username: "johndoe",
			Email:    "john@example.net",
		},
	}, users, "should return users response")

	mockRepoGetAll.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}
}

func testUserServiceRepoError(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{}, errors.New("repo_error"))

	users, err := service.GetAll()

	require.Nil(t, users, "should not return users")

	require.EqualError(t, err, "repo_error", "should return the repo error")

	mockRepoGetAll.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceGetAllSingle(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{
		{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		},
	}, nil)

	users, err := service.GetAll()

	require.NoError(t, err, "should not return an error")

	require.Equal(t, &[]UserResponseDto{
		{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		},
	}, users, "should return users response")

	mockRepoGetAll.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}
}

func testUserServiceGetAllZero(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{}, nil)

	users, err := service.GetAll()

	require.NoError(t, err, "should not return an error")

	require.Equal(t, &[]UserResponseDto{}, users, "should return users response")

	mockRepoGetAll.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}
}

func testUserServiceGetByAttr(t *testing.T, service UserService) {
	mockRepoGetByAttribute := mockRepo.
		On("GetByAttribute", "username", "janedoe").
		Return(&User{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		}, nil)

	user, err := service.GetByAttribute("username", "janedoe")

	require.NoError(t, err, "should not return error")
	require.Equal(t, &UserResponseDto{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, user, "should return matching user")

	mockRepoGetByAttribute.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceGetByAttrRepoError(t *testing.T, service UserService) {
	mockRepoGetByAttribute := mockRepo.
		On("GetByAttribute", "username", "janedoe").
		Return(&User{}, errors.New("repo_error"))

	user, err := service.GetByAttribute("username", "janedoe")

	require.EqualError(t, err, "repo_error", "should not return error")
	require.Nil(t, user, "should not return user")

	mockRepoGetByAttribute.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceUserExistsTrue(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.On("Exists", "janedoe").Return(true, nil)

	exists, err := service.Exists("janedoe")

	require.NoError(t, err, "should not return error")
	require.True(t, exists, "should return true")

	mockRepoUserExists.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceUserExistsFalse(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.On("Exists", "janedoe").Return(false, nil)

	exists, err := service.Exists("janedoe")

	require.NoError(t, err, "should not return error")
	require.False(t, exists, "should return false")

	mockRepoUserExists.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

}

func testUserServiceUserExistsError(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.
		On("Exists", "janedoe").
		Return(false, errors.New("repo_error"))

	exists, err := service.Exists("janedoe")

	require.EqualError(t, err, "repo_error", "should return error")
	require.False(t, exists, "should return false")

	mockRepoUserExists.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceDeleteById(t *testing.T, service UserService) {
	mockRepoDeleteById := mockRepo.On("DeleteById", uint(23)).Return(nil)

	err := service.DeleteById(uint(23))

	require.NoError(t, err, "should not return error")

	mockRepoDeleteById.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}

func testUserServiceDeleteByIdError(t *testing.T, service UserService) {
	mockRepoDeleteById := mockRepo.On("DeleteById", uint(23)).Return(errors.New("repo_error"))

	err := service.DeleteById(uint(23))

	require.EqualError(t, err, "repo_error", "should return error")

	mockRepoDeleteById.Unset()

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}
