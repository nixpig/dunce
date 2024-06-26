package user

import (
	"errors"
	"testing"

	"github.com/nixpig/dunce/pkg/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockRepo = new(MockUserRepo)
var mockCrypto = new(MockCrypto)

func TestUserService(t *testing.T) {
	scenarios := map[string]func(t *testing.T, service UserService){
		"get all users (success - multiple results)":                    testUserServiceGetAllMultiple,
		"get all users (success - zero results)":                        testUserServiceGetAllZero,
		"get all users (success - single result)":                       testUserServiceGetAllSingle,
		"get all users (error - repo)":                                  testUserServiceRepoError,
		"get by attribute (success)":                                    testUserServiceGetByAttr,
		"get by attribute (error - repo)":                               testUserServiceGetByAttrRepoError,
		"user exists (success - true)":                                  testUserServiceUserExistsTrue,
		"user exists (success - false)":                                 testUserServiceUserExistsFalse,
		"user exists (error - false)":                                   testUserServiceUserExistsError,
		"delete user by id (success)":                                   testUserServiceDeleteById,
		"delete user by id (error)":                                     testUserServiceDeleteByIdError,
		"create user (success)":                                         testUserServiceCreateUser,
		"create user (error - password hashing)":                        testUserServiceCreateUserPasswordHashingError,
		"create user (error - validation)":                              testUserServiceCreateUserValidationError,
		"create user (error - repo)":                                    testUserServiceCreateUserRepoError,
		"update user (success)":                                         testUserServiceUpdateUser,
		"update user (error - repo)":                                    testUserServiceUpdateUserRepoError,
		"update user (error - validation)":                              testUserServiceUpdateUserValidationError,
		"login with username and password (success)":                    testUserServiceLoginUsernamePassword,
		"login with username and password (error - repo)":               testUserServiceLoginUsernamePasswordRepoError,
		"login with username and password (error - incorrect password)": testUserServiceLoginUsernamePasswordIncorrectPassword,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			validator, err := validation.NewValidator()
			if err != nil {
				t.Fatal("unable to construct validator")
			}

			service := NewUserService(mockRepo, validator, mockCrypto)

			fn(t, service)
		})
	}
}

type MockCrypto struct {
	mock.Mock
}

func (mc MockCrypto) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := mc.Called(password, cost)

	return args.Get(0).([]byte), args.Error(1)
}

func (mc MockCrypto) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	args := mc.Called(hashedPassword, password)

	return args.Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

func (mu *MockUserRepo) Create(user *User) (*User, error) {
	args := mu.Called(user)

	return args.Get(0).(*User), args.Error(1)
}

func (mu *MockUserRepo) DeleteById(id uint) error {
	args := mu.Called(id)

	return args.Error(0)
}

func (mu *MockUserRepo) Exists(username string) (bool, error) {
	args := mu.Called(username)

	return args.Get(0).(bool), args.Error(1)
}

func (mu *MockUserRepo) GetAll() (*[]User, error) {
	args := mu.Called()

	return args.Get(0).(*[]User), args.Error(1)
}

func (mu *MockUserRepo) GetByAttribute(attr, value string) (*User, error) {
	args := mu.Called(attr, value)

	return args.Get(0).(*User), args.Error(1)
}

func (mu *MockUserRepo) Update(user *User) (*User, error) {
	args := mu.Called(user)

	return args.Get(0).(*User), args.Error(1)
}

func (mu *MockUserRepo) GetPasswordByUsername(username string) (string, error) {
	args := mu.Called(username)

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

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}

	mockRepoGetAll.Unset()
}

func testUserServiceRepoError(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{}, errors.New("repo_error"))

	users, err := service.GetAll()

	require.Nil(t, users, "should not return users")

	require.EqualError(t, err, "repo_error", "should return the repo error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetAll.Unset()
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

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}

	mockRepoGetAll.Unset()
}

func testUserServiceGetAllZero(t *testing.T, service UserService) {
	mockRepoGetAll := mockRepo.On("GetAll").Return(&[]User{}, nil)

	users, err := service.GetAll()

	require.NoError(t, err, "should not return an error")

	require.Equal(t, &[]UserResponseDto{}, users, "should return users response")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("expectations not met")
	}

	mockRepoGetAll.Unset()
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

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetByAttribute.Unset()
}

func testUserServiceGetByAttrRepoError(t *testing.T, service UserService) {
	mockRepoGetByAttribute := mockRepo.
		On("GetByAttribute", "username", "janedoe").
		Return(&User{}, errors.New("repo_error"))

	user, err := service.GetByAttribute("username", "janedoe")

	require.EqualError(t, err, "repo_error", "should not return error")
	require.Nil(t, user, "should not return user")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetByAttribute.Unset()
}

func testUserServiceUserExistsTrue(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.On("Exists", "janedoe").Return(true, nil)

	exists, err := service.Exists("janedoe")

	require.NoError(t, err, "should not return error")
	require.True(t, exists, "should return true")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoUserExists.Unset()
}

func testUserServiceUserExistsFalse(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.On("Exists", "janedoe").Return(false, nil)

	exists, err := service.Exists("janedoe")

	require.NoError(t, err, "should not return error")
	require.False(t, exists, "should return false")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoUserExists.Unset()
}

func testUserServiceUserExistsError(t *testing.T, service UserService) {
	mockRepoUserExists := mockRepo.
		On("Exists", "janedoe").
		Return(false, errors.New("repo_error"))

	exists, err := service.Exists("janedoe")

	require.EqualError(t, err, "repo_error", "should return error")
	require.False(t, exists, "should return false")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoUserExists.Unset()
}

func testUserServiceDeleteById(t *testing.T, service UserService) {
	mockRepoDeleteById := mockRepo.On("DeleteById", uint(23)).Return(nil)

	err := service.DeleteById(uint(23))

	require.NoError(t, err, "should not return error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoDeleteById.Unset()
}

func testUserServiceDeleteByIdError(t *testing.T, service UserService) {
	mockRepoDeleteById := mockRepo.On("DeleteById", uint(23)).Return(errors.New("repo_error"))

	err := service.DeleteById(uint(23))

	require.EqualError(t, err, "repo_error", "should return error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoDeleteById.Unset()
}

func testUserServiceCreateUser(t *testing.T, service UserService) {
	mockCryptoGenerateFromPassword := mockCrypto.
		On("GenerateFromPassword", []byte("foo"), 14).
		Return([]byte("hashed_password"), nil)

	mockRepoCreate := mockRepo.
		On("Create", &User{
			Username: "janedoe",
			Email:    "jane@example.org",
			Password: "hashed_password",
		}).
		Return(&User{
			Id:       23,
			Username: "janedoe",
			Email:    "jane@example.org",
		}, nil)

	createdUser, err := service.Create(&UserNewRequestDto{
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "foo",
	})

	require.NoError(t, err, "should not return error")
	require.Equal(t, &UserResponseDto{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, createdUser, "should return created user response")

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockCryptoGenerateFromPassword.Unset()
	mockRepoCreate.Unset()
}

func testUserServiceCreateUserPasswordHashingError(t *testing.T, service UserService) {
	mockCryptoGenerateFromPassword := mockCrypto.
		On("GenerateFromPassword", []byte("foo"), 14).
		Return([]byte(""), errors.New("password_error"))

	createdUser, err := service.Create(&UserNewRequestDto{
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "foo",
	})

	require.EqualError(t, err, "password_error", "should return error")
	require.Nil(t, createdUser, "should not return user")

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockCryptoGenerateFromPassword.Unset()
}

func testUserServiceCreateUserValidationError(t *testing.T, service UserService) {
	mockCryptoGenerateFromPassword := mockCrypto.
		On("GenerateFromPassword", []byte(""), 14).
		Return([]byte("hashed_password"), nil)

	createdUser, err := service.Create(&UserNewRequestDto{})

	require.Error(t, err, "should return error")
	require.Nil(t, createdUser, "should not return user")

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockCryptoGenerateFromPassword.Unset()
}

func testUserServiceCreateUserRepoError(t *testing.T, service UserService) {
	mockCryptoGenerateFromPassword := mockCrypto.
		On("GenerateFromPassword", []byte("foo"), 14).
		Return([]byte("hashed_password"), nil)

	mockRepoCreate := mockRepo.
		On("Create", &User{
			Username: "janedoe",
			Email:    "jane@example.org",
			Password: "hashed_password",
		}).
		Return(&User{}, errors.New("repo_error"))

	createdUser, err := service.Create(&UserNewRequestDto{
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "foo",
	})

	require.EqualError(t, err, "repo_error", "should return error")
	require.Nil(t, createdUser, "should not return user")

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockCryptoGenerateFromPassword.Unset()
	mockRepoCreate.Unset()
}

func testUserServiceUpdateUser(t *testing.T, service UserService) {
	mockRepoUpdate := mockRepo.On("Update", &User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "p4ssw0rd",
	}).Return(&User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, nil)

	updatedUser, err := service.Update(&User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "p4ssw0rd",
	})

	require.NoError(t, err, "should not return error")

	require.Equal(t, &UserResponseDto{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, updatedUser, "should return updated user response")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoUpdate.Unset()
}

func testUserServiceUpdateUserRepoError(t *testing.T, service UserService) {
	mockRepoUpdate := mockRepo.On("Update", &User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "p4ssw0rd",
	}).Return(&User{}, errors.New("repo_error"))

	updatedUser, err := service.Update(&User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "p4ssw0rd",
	})

	require.EqualError(t, err, "repo_error", "should return repo error")

	require.Nil(t, updatedUser, "should not return a user")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoUpdate.Unset()
}

func testUserServiceUpdateUserValidationError(t *testing.T, service UserService) {
	updatedUser, err := service.Update(&User{})

	require.Error(t, err, "should return repo error")

	require.Nil(t, updatedUser, "should not return a user")
}

func testUserServiceLoginUsernamePassword(t *testing.T, service UserService) {
	mockRepoGetPasswordByUserName := mockRepo.
		On("GetPasswordByUsername", "janedoe").
		Return("h4shedp4ssw0rd", nil)

	mockCryptoCompareHashAndPassword := mockCrypto.
		On("CompareHashAndPassword", []byte("h4shedp4ssw0rd"), []byte("p4ssw0rd")).
		Return(nil)

	err := service.LoginWithUsernamePassword("janedoe", "p4ssw0rd")

	require.NoError(t, err, "should not return error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetPasswordByUserName.Unset()
	mockCryptoCompareHashAndPassword.Unset()
}

func testUserServiceLoginUsernamePasswordRepoError(t *testing.T, service UserService) {
	mockRepoGetPasswordByUserName := mockRepo.
		On("GetPasswordByUsername", "janedoe").
		Return("", errors.New("repo_error"))

	err := service.LoginWithUsernamePassword("janedoe", "p4ssw0rd")

	require.EqualError(t, err, "repo_error", "should return repo error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetPasswordByUserName.Unset()
}

func testUserServiceLoginUsernamePasswordIncorrectPassword(t *testing.T, service UserService) {
	mockRepoGetPasswordByUserName := mockRepo.
		On("GetPasswordByUsername", "janedoe").
		Return("h4shedp4ssw0rd", nil)

	mockCryptoCompareHashAndPassword := mockCrypto.
		On("CompareHashAndPassword", []byte("h4shedp4ssw0rd"), []byte("p4ssw0rd")).
		Return(errors.New("incorrect_password"))

	err := service.LoginWithUsernamePassword("janedoe", "p4ssw0rd")

	require.EqualError(t, err, "incorrect_password", "should return crypto error")

	if res := mockRepo.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	if res := mockCrypto.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}

	mockRepoGetPasswordByUserName.Unset()
	mockCryptoCompareHashAndPassword.Unset()
}
