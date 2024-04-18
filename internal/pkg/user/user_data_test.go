package user_test

import (
	"regexp"
	"testing"

	"github.com/nixpig/dunce/internal/pkg/user"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestUserCreate(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, data user.UserData){
		"successfully create valid user": testUserCreate,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(*testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db connection pool")
			}

			defer mock.Close()

			data := user.NewUserData(mock)

			fn(t, mock, data)
		})
	}

}

func testUserCreate(t *testing.T, mock pgxmock.PgxPoolIface, data user.UserData) {
	query := `insert into users_ (username_, email_, link_, role_, password_) values($1, $2, $3, $4, $5) returning id_, username_, email_, link_, role_`

	row := pgxmock.
		NewRows([]string{"id_", "username_", "email_", "link_", "role_"}).
		AddRow(23, "janedoe", "jane@example.org", "https://t.com/jane", "author")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe", "jane@example.org", "https://t.com/jane", "author", "p4ssw0rd").
		WillReturnRows(row)

	createdUser, err := data.Save(user.UserNew{
		Username: "janedoe",
		Email:    "jane@example.org",
		Link:     "https://t.com/jane",
		Role:     "author",
		Password: "p4ssw0rd",
	})

	require.Nil(t, err, "should not return error")

	require.Equal(t, &user.User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
		Link:     "https://t.com/jane",
		Role:     "author",
	}, createdUser, "should return created user response")
}
