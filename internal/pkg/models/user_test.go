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

func TestCreateUser(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T, mock pgxmock.PgxPoolIface){
		"inserts valid user details":    testInsertValidUser,
		"reject duplicate user details": testRejectDuplicateUser,
		"reject invalid user details":   testRejectInvalidUser,
		"reject empty user details":     testRejectEmptyUser,
	} {
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

// func TestGetUsers(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatalf("unable to create mock database pool")
// 	}
//
// 	defer mock.Close()
//
// 	models.BuildQueries(mock)
//
// 	mockRows := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
// 		AddRow(1, "user1", "test_one@example.com", "http://example.com", models.AdminRole).
// 		AddRow(2, "user2", "test_two@example.org", "https://example.org", models.AuthorRole)
//
// 	mock.ExpectQuery(`select id_, username_, email_, link_, role_ from user_`).WillReturnRows(mockRows)
//
// 	users, err := models.Query.User.GetAll()
// 	if err != nil {
// 		t.Fatalf("unable to get users: %v", err)
// 	}
//
// 	assert.Equal(t, &[]models.User{
// 		{Id: 1, UserData: models.UserData{Username: "user1", Email: "test_one@example.com", Link: "http://example.com", Role: models.AdminRole}},
// 		{Id: 2, UserData: models.UserData{Username: "user2", Email: "test_two@example.org", Link: "https://example.org", Role: models.AuthorRole}},
// 	}, users)
// }
//
// func TestGetUsersByRole(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatalf("unable to create mock database: %v", err)
// 	}
//
// 	defer mock.Close()
//
// 	models.BuildQueries(mock)
//
// 	mockRowsAdmin := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
// 		AddRow(1, "admin1", "admin1@example.com", "https://admin1-example.org", models.AdminRole).
// 		AddRow(6, "admin2", "admin2@example.com", "https://admin2-example.org", models.AdminRole)
//
// 	mock.
// 		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
// 		WithArgs(models.AdminRole).
// 		WillReturnRows(mockRowsAdmin)
//
// 	adminUsers, err := models.Query.User.GetByRole(models.AdminRole)
// 	if err != nil {
// 		t.Fatalf("unable to get users with Admin role: %v", err)
// 	}
//
// 	assert.Equal(t, &[]models.User{
// 		{Id: 1, UserData: models.UserData{Username: "admin1", Email: "admin1@example.com", Link: "https://admin1-example.org", Role: models.AdminRole}},
// 		{Id: 6, UserData: models.UserData{Username: "admin2", Email: "admin2@example.com", Link: "https://admin2-example.org", Role: models.AdminRole}},
// 	}, adminUsers)
//
// 	mockRowsAuthor := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
// 		AddRow(2, "author1", "author1@example.com", "https://author1-example.org", models.AuthorRole).
// 		AddRow(4, "author2", "author2@example.com", "https://author2-example.org", models.AuthorRole).
// 		AddRow(5, "author3", "author3@example.com", "https://author3-example.org", models.AuthorRole)
//
// 	mock.
// 		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
// 		WithArgs(models.AuthorRole).
// 		WillReturnRows(mockRowsAuthor)
//
// 	authorUsers, err := models.Query.User.GetByRole(models.AuthorRole)
// 	if err != nil {
// 		t.Fatalf("unable to get users with Author role: %v", err)
// 	}
//
// 	assert.Equal(t, &[]models.User{
// 		{Id: 2, UserData: models.UserData{Username: "author1", Email: "author1@example.com", Link: "https://author1-example.org", Role: models.AuthorRole}},
// 		{Id: 4, UserData: models.UserData{Username: "author2", Email: "author2@example.com", Link: "https://author2-example.org", Role: models.AuthorRole}},
// 		{Id: 5, UserData: models.UserData{Username: "author3", Email: "author3@example.com", Link: "https://author3-example.org", Role: models.AuthorRole}},
// 	}, authorUsers)
//
// 	mockRowsReader := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
// 		AddRow(3, "reader1", "reader1@example.com", "https://reader1-example.org", models.ReaderRole).
// 		AddRow(7, "reader2", "reader2@example.com", "https://reader2-example.org", models.ReaderRole)
//
// 	mock.
// 		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
// 		WithArgs(models.ReaderRole).
// 		WillReturnRows(mockRowsReader)
//
// 	readerUsers, err := models.Query.User.GetByRole(models.ReaderRole)
// 	if err != nil {
// 		t.Fatalf("unable to get users with Reader role: %v", err)
// 	}
//
// 	assert.Equal(t, &[]models.User{
// 		{Id: 3, UserData: models.UserData{Username: "reader1", Email: "reader1@example.com", Link: "https://reader1-example.org", Role: models.ReaderRole}},
// 		{Id: 7, UserData: models.UserData{Username: "reader2", Email: "reader2@example.com", Link: "https://reader2-example.org", Role: models.ReaderRole}},
// 	}, readerUsers)
//
// 	mockRowsEmpty := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"})
//
// 	mock.
// 		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
// 		WithArgs("UNKNOWN_ROLE").
// 		WillReturnRows(mockRowsEmpty)
//
// 	unknown, err := models.Query.User.GetByRole("UNKNOWN")
// 	assert.NotNil(t, err)
// 	assert.Nil(t, unknown)
//
// }
//
// func TestGetUserById(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatalf("unable to create mock database: %v", err)
// 	}
//
// 	defer mock.Close()
//
// 	models.BuildQueries(mock)
//
// 	mockRow := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
// 		AddRow(23, "test_user", "test@example.com", "http://example.com", models.AuthorRole)
//
// 	mock.ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where id_ = $1`)).
// 		WithArgs(23).
// 		WillReturnRows(mockRow)
//
// 	user23, err := models.Query.User.GetById(23)
// 	if err != nil {
// 		t.Fatalf("unable to get user 23: %v", err)
// 	}
//
// 	assert.Equal(t, &models.User{
// 		Id: 23,
// 		UserData: models.UserData{Username: "test_user",
// 			Email: "test@example.com",
// 			Link:  "http://example.com",
// 			Role:  models.AuthorRole,
// 		},
// 	}, user23)
//
// 	mock.ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where id_ = $1`)).
// 		WithArgs(7).
// 		WillReturnRows(mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}))
//
// 	user7, err := models.Query.User.GetById(7)
// 	assert.NotNil(t, err)
// 	assert.Nil(t, user7)
// }
