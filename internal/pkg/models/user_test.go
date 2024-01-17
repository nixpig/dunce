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
		t.Fatalf("unable to get mock database")
	}

	defer mock.Close()

	api := models.Api{DB: mock}

	newUser := models.NewUser{
		Username: "some username",
		Email:    "some email",
		Link:     "some link",
		Role:     models.AdminRole,
		Password: "some password",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`insert into user_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`)).
		WithArgs(newUser.Username, newUser.Email, newUser.Link, newUser.Role, newUser.Password).
		WillReturnRows(mock.NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).AddRow(23, newUser.Username, newUser.Email, newUser.Link, newUser.Role))

	createdUser, err := api.CreateUser(&newUser)
	if err != nil {
		t.Errorf("%v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%v", err)
	}

	assert.Equal(t, &models.User{
		Id:       23,
		Username: newUser.Username,
		Email:    newUser.Email,
		Link:     newUser.Link,
		Role:     newUser.Role,
	}, createdUser)

}
