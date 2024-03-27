package models_test

import (
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AnyPassword struct{}

func (a AnyPassword) Match(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func TestDeleteUserById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"delete existing user by ID":     testDeleteExistingUserById,
		"delete non-existent user by ID": testDeleteNonExistentUserById,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetUsersByRole(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get users by role - multiple results": getUsersByRoleMultipleResults,
		"get users by role - single result":    getUsersByRoleSingleResult,
		"get users by role - no results":       getUsersByRoleNoResults,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create database connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get user with valid email":   testGetUserWithValidEmail,
		"get user with invalid email": testGetUserWithInvalidEmail,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}

}

func TestGetUserByUsername(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"get user with valid username":   testGetUserWithValidUsername,
		"get user with invalid username": testGetUserWithInvalidUsername,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("unable to create db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestCreateUser(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"inserts valid user details":    testInsertValidUser,
		"reject duplicate user details": testRejectDuplicateUser,
		"reject invalid user details":   testRejectInvalidUser,
		"reject empty user details":     testRejectEmptyUser,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("unable to create mock database pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"update valid user details":          testUpdateValidUser,
		"reject update invalid user details": testUpdateInvalidUser,
		"reject update non-existent user":    testUpdateNonExistentUser,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("failed to create mock connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"test get all users with multiple results": testGetAllUsersMultipleResults,
		"test get all users with single result":    testGetAllUsersSingleResult,
		"test get all users with no results":       testGetAllUsersNoResults,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("unable to create mock db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func TestGetUserById(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"test get user with valid id":   testGetUserWithValidId,
		"test get user with invalid id": testGetUserWithInvalidId,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("unable to create mock db connection pool")
			}

			defer mock.Close()

			models.BuildQueries(mock)

			fn(t, mock)
		})
	}
}

func testGetUserWithValidId(t *testing.T, mock pgxmock.PgxPoolIface) {
	mockResult := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "user", "user@example.org", "https://t.com/u", models.AdminRole)

	selectQuery := "select id_, username_, email_, link_, role_ from users_ where id_ = $1"

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs(23).WillReturnRows(mockResult)

	user, err := models.Query.User.GetById(23)
	require.Nil(t, err, "expected not to error")

	require.Equal(t, &models.User{
		Id: 23,
		UserData: models.UserData{
			Username: "user",
			Email:    "user@example.org",
			Link:     "https://t.com/u",
			Role:     models.AdminRole,
		},
	}, user)
}

func testGetUserWithInvalidId(t *testing.T, mock pgxmock.PgxPoolIface) {
	emptyResult := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	selectQuery := "select id_, username_, email_, link_, role_ from users_ where id_ = $1"

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs(23).WillReturnRows(emptyResult)

	user, err := models.Query.User.GetById(23)

	require.Nil(t, user, "user should be empty")
	require.NotNil(t, err, "should return an error")
}

func testGetAllUsersMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	multipleResults := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRows([][]any{
			{
				23, "firstuser", "first@example.org", "https://t.com/f", models.AuthorRole,
			},
			{
				42, "seconduser", "second@example.org", "https://t.com/s", models.ReaderRole,
			},
		}...)

	selectQuery := `select id_, username_, email_, link_, role_ from users_`

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(multipleResults)

	users, err := models.Query.User.GetAll()
	require.Nil(t, err)
	require.Equal(t, &[]models.User{
		{
			Id: 23,
			UserData: models.UserData{
				Username: "firstuser",
				Email:    "first@example.org",
				Link:     "https://t.com/f",
				Role:     models.AuthorRole,
			},
		},
		{
			Id: 42,
			UserData: models.UserData{
				Username: "seconduser",
				Email:    "second@example.org",
				Link:     "https://t.com/s",
				Role:     models.ReaderRole,
			},
		},
	}, users)
}

func testGetAllUsersSingleResult(t *testing.T, mock pgxmock.PgxPoolIface) {
	singleResult := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "firstuser", "first@example.org", "https://t.com/f", models.AdminRole)

	selectQuery := `select id_, username_, email_, link_, role_ from users_`

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(singleResult)

	users, err := models.Query.User.GetAll()
	require.Nil(t, err)
	require.Equal(t, &[]models.User{
		{
			Id: 23,
			UserData: models.UserData{
				Username: "firstuser",
				Email:    "first@example.org",
				Link:     "https://t.com/f",
				Role:     models.AdminRole,
			},
		},
	}, users)
}

func testGetAllUsersNoResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	emptyResult := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	selectQuery := `select id_, username_, email_, link_, role_ from users_`

	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(emptyResult)

	users, err := models.Query.User.GetAll()
	require.Nil(t, err)
	require.Equal(t, users, &[]models.User{})
}

func testUpdateNonExistentUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	updatedUser := models.UserData{
		Username: "differentname",
		Email:    "different@somewhere.com",
		Link:     "https://twitter.com",
		Role:     models.AdminRole,
	}

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	mockUpdateQuery := `update users_ set username_ = $2, email_ = $3, link_ = $4, role_ = $5 where id_ = $1 returning id_, username_, email_, link_, role_`
	mock.
		ExpectQuery(regexp.QuoteMeta(mockUpdateQuery)).
		WithArgs(23, &updatedUser.Username, &updatedUser.Email, &updatedUser.Link, &updatedUser.Role).
		WillReturnRows(mockRows)

	user, err := models.Query.User.UpdateById(23, &updatedUser)
	require.Error(t, err)
	require.Equal(t, models.UserError{"User does not exist"}, err)
	require.Nil(t, user)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func testUpdateInvalidUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	invalidUser1 := models.UserData{
		Username: "u",     // username too short
		Email:    "foo",   // invalid email address
		Link:     "bar",   // invalid URL
		Role:     "admin", // not a valid role
	}

	createdUser1, err := models.Query.User.UpdateById(23, &invalidUser1)
	require.Nil(t, createdUser1)
	require.Equal(t, models.UserError{
		"Username field requires a min length of 5; length of value provided is 1",
		"Email field requires an email but received foo",
		"Link field requires a URL but received bar",
	}, err)

	invalidUser2 := models.UserData{
		Username: "username is way too long to be valid", // too long username
		Email:    "test_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_real@somewhere.com",
		Link:     "https://link_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likinglink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likingslink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likingsslink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likings.com/some_really_long_path_under_the_url_link",
		Role:     "author",
	}

	createdUser2, err := models.Query.User.UpdateById(23, &invalidUser2)
	require.Nil(t, createdUser2)
	require.Equal(t, models.UserError{
		"Username field has a max length of 16; length of value provided is 36",
		"Email field has a max length of 100; length of value provided is 262",
		"Link field has a max length of 255; length of value provided is 361",
	}, err)

	invalidPassword := "foobar"
	validUser := models.UserData{
		Username: "username",
		Email:    "somebody@somewhere.com",
		Link:     "https://somewhere.com/somebody",
		Role:     models.AdminRole,
	}

	createdInvalidUser, err := models.Query.User.Create(&validUser, invalidPassword)
	require.Nil(t, createdInvalidUser)
	require.Equal(t, models.UserError{"Password must be longer than 7 characters"}, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met")
	}
}

func testUpdateValidUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	updatedUser := models.UserData{
		Username: "differentname",
		Email:    "different@somewhere.com",
		Link:     "https://twitter.com",
		Role:     models.AdminRole,
	}

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, updatedUser.Username, updatedUser.Email, updatedUser.Link, updatedUser.Role)

	mockUpdateQuery := `update users_ set username_ = $2, email_ = $3, link_ = $4, role_ = $5 where id_ = $1 returning id_, username_, email_, link_, role_`
	mock.
		ExpectQuery(regexp.QuoteMeta(mockUpdateQuery)).
		WithArgs(23, &updatedUser.Username, &updatedUser.Email, &updatedUser.Link, &updatedUser.Role).
		WillReturnRows(mockRows)

	user, err := models.Query.User.UpdateById(23, &updatedUser)
	require.NoError(t, err)
	require.Equal(t, &models.User{
		Id: 23,
		UserData: models.UserData{
			Username: updatedUser.Username,
			Email:    "different@somewhere.com",
			Link:     "https://twitter.com",
			Role:     models.AdminRole,
		},
	}, user)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func testRejectDuplicateUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	existingUser := models.UserData{
		Username: "somebody",
		Email:    "sombody@example.org",
		Link:     "",
		Role:     models.AuthorRole,
	}

	dupeCheckQuery := `select count(id_) from users_ where username_ = $1 or email_ = $2`

	dupeUsernameRows := mock.
		NewRows([]string{"count"}).
		AddRow(1)

	mock.
		ExpectQuery(regexp.QuoteMeta(dupeCheckQuery)).
		WithArgs(existingUser.Username, "not_a_duplicate@example.net").
		WillReturnRows(dupeUsernameRows)

	dupeUsernameUser, err := models.Query.User.Create(&models.UserData{
		Username: existingUser.Username,
		Email:    "not_a_duplicate@example.net",
		Link:     "",
		Role:     models.ReaderRole,
	}, "somepassword")

	require.Nil(t, dupeUsernameUser)
	require.Equal(t, models.UserError{"User already exists"}, err)

	dupeEmailRows := mock.
		NewRows([]string{"count"}).
		AddRow(1)

	mock.
		ExpectQuery(regexp.QuoteMeta(dupeCheckQuery)).
		WithArgs("uniqueusername", existingUser.Email).
		WillReturnRows(dupeEmailRows)

	dupeEmailUser, err := models.Query.User.Create(&models.UserData{
		Username: "uniqueusername",
		Email:    existingUser.Email,
		Link:     "",
		Role:     models.ReaderRole,
	}, "somepassword")

	require.Nil(t, dupeEmailUser)
	require.Equal(t, models.UserError{"User already exists"}, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met")
	}
}

func testRejectInvalidUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	invalidUser1 := models.UserData{
		Username: "u",     // username too short
		Email:    "foo",   // invalid email address
		Link:     "bar",   // invalid URL
		Role:     "admin", // not a valid role
	}

	createdUser1, err := models.Query.User.Create(&invalidUser1, "some password")
	require.Nil(t, createdUser1)
	require.Equal(t, models.UserError{
		"Username field requires a min length of 5; length of value provided is 1",
		"Email field requires an email but received foo",
		"Link field requires a URL but received bar",
	}, err)

	invalidUser2 := models.UserData{
		Username: "username is way too long to be valid", // too long username
		Email:    "test_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_realtest_name_for_an_email_address_that_is_way_too_long_to_be_real@somewhere.com",
		Link:     "https://link_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likinglink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likingslink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likingsslink_that_may_well_be_valid_in_structure_but_that_is_too_long_for_our_likings.com/some_really_long_path_under_the_url_link",
		Role:     "author",
	}

	createdUser2, err := models.Query.User.Create(&invalidUser2, "some password")
	require.Nil(t, createdUser2)
	require.Equal(t, models.UserError{
		"Username field has a max length of 16; length of value provided is 36",
		"Email field has a max length of 100; length of value provided is 262",
		"Link field has a max length of 255; length of value provided is 361",
	}, err)

	invalidPassword := "foobar"
	validUser := models.UserData{
		Username: "username",
		Email:    "somebody@somewhere.com",
		Link:     "https://somewhere.com/somebody",
		Role:     models.AdminRole,
	}

	createdInvalidUser, err := models.Query.User.Create(&validUser, invalidPassword)
	require.Nil(t, createdInvalidUser)
	require.Equal(t, models.UserError{"Password must be longer than 7 characters"}, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met")
	}
}

func testRejectEmptyUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	newUser := models.UserData{
		Username: "",
		Email:    "",
		Link:     "",
	}

	password := ""

	createdUser, err := models.Query.User.Create(&newUser, password)
	require.Nil(t, createdUser)
	require.Equal(t, models.UserError{
		"Username field is required but is empty",
		"Email field is required but is empty",
		"Role field is required but is empty",
	}, err)
}

func testInsertValidUser(t *testing.T, mock pgxmock.PgxPoolIface) {
	newUser := models.UserData{
		Username: "username",
		Email:    "somebody@somewhere.com",
		Link:     "https://somewhere.com/somebody",
		Role:     models.AdminRole,
	}

	password := "some password"

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, newUser.Username, newUser.Email, newUser.Link, newUser.Role)

	dupeCheckQuery := `select count(id_) from users_ where username_ = $1 or email_ = $2 `
	expectQuery := `insert into users_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	mock.
		ExpectQuery(regexp.QuoteMeta(dupeCheckQuery)).
		WithArgs(newUser.Username, newUser.Email).
		WillReturnRows(mock.NewRows([]string{"id_"}))

	mock.
		ExpectQuery(regexp.QuoteMeta(expectQuery)).
		WithArgs(
			newUser.Username,
			newUser.Email,
			newUser.Link,
			newUser.Role,
			AnyPassword{},
		).
		WillReturnRows(mockRows)

	createdUser, err := models.Query.User.Create(&newUser, password)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}

	assert.Equal(t, &models.User{
		Id: 23,
		UserData: models.UserData{
			Username: newUser.Username,
			Email:    newUser.Email,
			Link:     newUser.Link,
			Role:     newUser.Role,
		},
	}, createdUser)
}

func testGetUserWithValidUsername(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where username_ = $1"

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "someuser", "somebody@example.com", "https://g.com/user", models.ReaderRole)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("someuser").WillReturnRows(mockRow)

	user, err := models.Query.User.GetByUsername("someuser")
	require.Nil(t, err, "should not error")
	require.Equal(t, &models.User{
		Id: 23,
		UserData: models.UserData{
			Username: "someuser",
			Email:    "somebody@example.com",
			Link:     "https://g.com/user",
			Role:     models.ReaderRole,
		},
	}, user)
}

func testGetUserWithInvalidUsername(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where username_ = $1"

	mockEmptyRows := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockEmptyRows)

	user, err := models.Query.User.GetByUsername("foobar")
	require.Nil(t, user, "no user should be returned")
	require.NotNil(t, err, "error should be returned")
}

func testGetUserWithValidEmail(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where email_ = $1"

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "someuser", "somebody@example.com", "https://g.com/user", models.ReaderRole)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("somebody@example.com").WillReturnRows(mockRow)

	user, err := models.Query.User.GetByEmail("somebody@example.com")

	require.Nil(t, err, "should not error")
	require.Equal(t, &models.User{
		Id: 23,
		UserData: models.UserData{
			Username: "someuser",
			Email:    "somebody@example.com",
			Link:     "https://g.com/user",
			Role:     models.ReaderRole,
		},
	}, user)
}

func testGetUserWithInvalidEmail(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where email_ = $1"

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("somebody@example.com").WillReturnRows(mockRow)

	user, err := models.Query.User.GetByEmail("somebody@example.com")
	require.NotNil(t, err, "should error")
	require.Nil(t, user, "should be nil")

}

func getUsersByRoleMultipleResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where role_ = $1"

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRows([][]any{
			{
				23, "userone", "userone@example.org", "https://t.com/1", models.ReaderRole,
			},
			{
				52, "usertwo", "usertwo@example.com", "https://t.com/2", models.ReaderRole,
			},
		}...)

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(models.ReaderRole).
		WillReturnRows(mockRows)

	users, err := models.Query.User.GetByRole(models.ReaderRole)
	require.Nil(t, err, "should not error")
	require.Equal(t, &[]models.User{
		{
			Id: 23,
			UserData: models.UserData{
				Username: "userone",
				Email:    "userone@example.org",
				Link:     "https://t.com/1",
				Role:     models.ReaderRole,
			},
		},
		{
			Id: 52,
			UserData: models.UserData{
				Username: "usertwo",
				Email:    "usertwo@example.com",
				Link:     "https://t.com/2",
				Role:     models.ReaderRole,
			},
		},
	}, users, "should return users by role")
}

func getUsersByRoleSingleResult(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where role_ = $1"

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRows([][]any{
			{

				23, "userone", "userone@example.org", "https://t.com/1", models.ReaderRole,
			},
		}...)

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(models.ReaderRole).
		WillReturnRows(mockRow)

	users, err := models.Query.User.GetByRole(models.ReaderRole)
	require.Nil(t, err, "should not error")
	require.Equal(t, &[]models.User{
		{
			Id: 23,
			UserData: models.UserData{
				Username: "userone",
				Email:    "userone@example.org",
				Link:     "https://t.com/1",
				Role:     models.ReaderRole,
			},
		},
	}, users, "should return users by role")
}

func getUsersByRoleNoResults(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := "select id_, username_, email_, link_, role_ from users_ where role_ = $1"

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(models.ReaderRole).
		WillReturnRows(mockRow)

	users, err := models.Query.User.GetByRole(models.ReaderRole)
	require.Nil(t, err, "should not be an error when no results")
	require.Empty(t, users, "should not return any users")
}

func testDeleteExistingUserById(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `delete from users_ where id_ = $1`

	mockResult := pgxmock.NewResult("delete", 1)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(23).WillReturnResult(mockResult)

	err := models.Query.User.DeleteById(23)

	require.Nil(t, err, "should not error out")
}

func testDeleteNonExistentUserById(t *testing.T, mock pgxmock.PgxPoolIface) {
	query := `delete from users_ where id_ = $1`

	mockResult := pgxmock.NewResult("delete", 0)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(42).WillReturnResult(mockResult)

	err := models.Query.User.DeleteById(42)

	require.NotNil(t, err, "should error out when no users deleted")
}
