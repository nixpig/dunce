package models_test

import (
	"regexp"
	"testing"

	"github.com/nixpig/bloggor/internal/pkg/models"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unable to create mock database pool")
	}

	defer mock.Close()

	models.BuildQueries(mock)

	newUser := models.NewUserData{
		Username: "some username",
		Email:    "some email",
		Link:     "some link",
		Role:     models.AdminRole,
		Password: "some password",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`)).
		WithArgs(newUser.Username, newUser.Email, newUser.Link, newUser.Role, newUser.Password).
		WillReturnRows(mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).AddRow(23, newUser.Username, newUser.Email, newUser.Link, newUser.Role))

	createdUser, err := models.Query.User.Create(&newUser)
	if err != nil {
		t.Errorf("%v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%v", err)
	}

	assert.Equal(t, &models.UserData{
		Id:       23,
		Username: newUser.Username,
		Email:    newUser.Email,
		Link:     newUser.Link,
		Role:     newUser.Role,
	}, createdUser)
}

func TestGetUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unable to create mock database pool")
	}

	defer mock.Close()

	models.BuildQueries(mock)

	mockRows := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(1, "user1", "test_one@example.com", "http://example.com", models.AdminRole).
		AddRow(2, "user2", "test_two@example.org", "https://example.org", models.AuthorRole)

	mock.ExpectQuery(`select id_, username_, email_, link_, role_ from user_`).WillReturnRows(mockRows)

	users, err := models.Query.User.GetAll()
	if err != nil {
		t.Fatalf("unable to get users: %v", err)
	}

	assert.Equal(t, &[]models.UserData{
		{Id: 1, Username: "user1", Email: "test_one@example.com", Link: "http://example.com", Role: models.AdminRole},
		{Id: 2, Username: "user2", Email: "test_two@example.org", Link: "https://example.org", Role: models.AuthorRole},
	}, users)
}

func TestGetUsersByRole(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unable to create mock database: %v", err)
	}

	defer mock.Close()

	models.BuildQueries(mock)

	mockRowsAdmin := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(1, "admin1", "admin1@example.com", "https://admin1-example.org", models.AdminRole).
		AddRow(6, "admin2", "admin2@example.com", "https://admin2-example.org", models.AdminRole)

	mock.
		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
		WithArgs(models.AdminRole).
		WillReturnRows(mockRowsAdmin)

	adminUsers, err := models.Query.User.GetByRole(models.AdminRole)
	if err != nil {
		t.Fatalf("unable to get users with Admin role: %v", err)
	}

	assert.Equal(t, &[]models.UserData{
		{Id: 1, Username: "admin1", Email: "admin1@example.com", Link: "https://admin1-example.org", Role: models.AdminRole},
		{Id: 6, Username: "admin2", Email: "admin2@example.com", Link: "https://admin2-example.org", Role: models.AdminRole},
	}, adminUsers)

	mockRowsAuthor := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(2, "author1", "author1@example.com", "https://author1-example.org", models.AuthorRole).
		AddRow(4, "author2", "author2@example.com", "https://author2-example.org", models.AuthorRole).
		AddRow(5, "author3", "author3@example.com", "https://author3-example.org", models.AuthorRole)

	mock.
		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
		WithArgs(models.AuthorRole).
		WillReturnRows(mockRowsAuthor)

	authorUsers, err := models.Query.User.GetByRole(models.AuthorRole)
	if err != nil {
		t.Fatalf("unable to get users with Author role: %v", err)
	}

	assert.Equal(t, &[]models.UserData{
		{Id: 2, Username: "author1", Email: "author1@example.com", Link: "https://author1-example.org", Role: models.AuthorRole},
		{Id: 4, Username: "author2", Email: "author2@example.com", Link: "https://author2-example.org", Role: models.AuthorRole},
		{Id: 5, Username: "author3", Email: "author3@example.com", Link: "https://author3-example.org", Role: models.AuthorRole},
	}, authorUsers)

	mockRowsReader := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(3, "reader1", "reader1@example.com", "https://reader1-example.org", models.ReaderRole).
		AddRow(7, "reader2", "reader2@example.com", "https://reader2-example.org", models.ReaderRole)

	mock.
		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
		WithArgs(models.ReaderRole).
		WillReturnRows(mockRowsReader)

	readerUsers, err := models.Query.User.GetByRole(models.ReaderRole)
	if err != nil {
		t.Fatalf("unable to get users with Reader role: %v", err)
	}

	assert.Equal(t, &[]models.UserData{
		{Id: 3, Username: "reader1", Email: "reader1@example.com", Link: "https://reader1-example.org", Role: models.ReaderRole},
		{Id: 7, Username: "reader2", Email: "reader2@example.com", Link: "https://reader2-example.org", Role: models.ReaderRole},
	}, readerUsers)

	mockRowsEmpty := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"})

	mock.
		ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where role_ = $1`)).
		WithArgs("UNKNOWN_ROLE").
		WillReturnRows(mockRowsEmpty)

	unknown, err := models.Query.User.GetByRole("UNKNOWN")
	assert.NotNil(t, err)
	assert.Nil(t, unknown)

}

func TestGetUserById(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unable to create mock database: %v", err)
	}

	defer mock.Close()

	models.BuildQueries(mock)

	mockRow := mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "test_user", "test@example.com", "http://example.com", models.AuthorRole)

	mock.ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where id_ = $1`)).
		WithArgs(23).
		WillReturnRows(mockRow)

	user23, err := models.Query.User.GetById(23)
	if err != nil {
		t.Fatalf("unable to get user 23: %v", err)
	}

	assert.Equal(t, &models.UserData{
		Id:       23,
		Username: "test_user",
		Email:    "test@example.com",
		Link:     "http://example.com",
		Role:     models.AuthorRole,
	}, user23)

	mock.ExpectQuery(regexp.QuoteMeta(`select id_, username_, email_, link_, role_ from user_ where id_ = $1`)).
		WithArgs(7).
		WillReturnRows(mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}))

	user7, err := models.Query.User.GetById(7)
	assert.NotNil(t, err)
	assert.Nil(t, user7)
}
