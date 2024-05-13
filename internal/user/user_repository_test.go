package user

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestUserRepo(t *testing.T) {
	scenarios := map[string]func(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository){
		"create user (success)":                              testUserRepoCreate,
		"create user (error)":                                testUserRepoCreateError,
		"delete user (success)":                              testUserRepoDelete,
		"delete user (error - db error)":                     testUserRepoDeleteDbError,
		"delete user (error - zero rows)":                    testUserRepoDeleteNoRows,
		"user exists (true)":                                 testUserRepoUserExists,
		"user exists (false)":                                testUserRepoUserNotExists,
		"user exists (false - db error)":                     testUserRepoExistsDbError,
		"get all users (success)":                            testUserRepoGetAll,
		"get all users (db error)":                           testUserRepoGetAllDbError,
		"get all users (scan error)":                         testUserRepoGetAllScanError,
		"get user by attribute - username (success)":         testUserRepoGetByUsername,
		"get user by attribute - unknown (error)":            testUserRepoGetByUnknownAttribute,
		"get user by attribute - username (error - db scan)": testUserRepoGetByAttributeDbError,
		"get password by username (success)":                 testUserRepoGetPasswordByUsername,
		"get password by username (error - db error)":        testUserRepoGetPasswordByUsernameDbError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			db, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal("unable to create mock db pool")
			}

			defer db.Close()

			repo := NewUserPostgresRepository(db)

			fn(t, db, repo)
		})
	}
}

func testUserRepoCreate(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `insert into users_ (username_, email_, password_) values ($1, $2, $3) returning id_, username_, email_`

	mockRow := mock.
		NewRows([]string{"id_", "username_", "email_"}).
		AddRow(uint(23), "janedoe", "jane@example.org")

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe", "jane@example.org", "p4ssw0rd").
		WillReturnRows(mockRow)

	createdUser, err := repo.Create(&User{
		Username: "janedoe",
		Email:    "jane@example.org",
		Password: "p4ssw0rd",
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations were not met", err)
	}

	require.Nil(t, err, "should not return error")
	require.Equal(t, &User{
		Id:       uint(23),
		Username: "janedoe",
		Email:    "jane@example.org",
	}, createdUser, "should return created user details")

}

func testUserRepoCreateError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `insert into users_ (username_, email_, password_) values ($1, $2, $3) returning id_, username_, email_`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe", "jane@example.com", "p4ssw0rd").
		WillReturnError(errors.New("db_error"))

	user, err := repo.Create(&User{
		Username: "janedoe",
		Email:    "jane@example.com",
		Password: "p4ssw0rd",
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	require.Nil(t, user, "should not return user details")
	require.EqualError(t, err, "db_error", "should return db error")

	mock.Reset()

}

func testUserRepoDelete(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `delete from users_ where id_ = $1`

	var id uint = 23

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("delete", 1))

	err := repo.DeleteById(id)

	require.Nil(t, err, "should not return error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoDeleteDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `delete from users_ where id_ = $1`

	var id uint = 69

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnError(errors.New("db_error"))

	err := repo.DeleteById(id)

	require.EqualError(t, err, "db_error", "should return db error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoDeleteNoRows(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `delete from users_ where id_ = $1`

	var id uint = 69

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("delete", 0))

	err := repo.DeleteById(id)

	require.Error(t, err, "should return zero rows error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()

}

func testUserRepoUserExists(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select exists(select true from users_ where username_ = $1)`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe").
		WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.Exists("janedoe")

	require.Nil(t, err, "should not return error")
	require.True(t, exists, "should return true")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoUserNotExists(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select exists(select true from users_ where username_ = $1)`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe").
		WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.Exists("janedoe")

	require.Nil(t, err, "should not return error")
	require.False(t, exists, "should return false")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoExistsDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select exists(select true from users_ where username_ = $1)`

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("janedoe").
		WillReturnError(errors.New("db_error"))

	exists, err := repo.Exists("janedoe")

	require.Error(t, err, "should return error")
	require.False(t, exists, "should return false")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()

}

func testUserRepoGetAll(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select id_, username_, email_ from users_`

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_"}).
		AddRow(uint(23), "janedoe", "jane@example.org").
		AddRow(uint(42), "johndoe", "john@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	users, err := repo.GetAll()

	require.Nil(t, err, "should not return error")

	require.Equal(t, &[]User{
		{
			Id:       uint(23),
			Username: "janedoe",
			Email:    "jane@example.org",
		},
		{
			Id:       uint(42),
			Username: "johndoe",
			Email:    "john@example.com",
		},
	}, users, "should return user list")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoGetAllDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select id_, username_, email_ from users_`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("db_error"))

	users, err := repo.GetAll()

	require.Nil(t, users, "should not return users")

	require.EqualError(t, err, "db_error", "should return db error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()

}

func testUserRepoGetAllScanError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select id_, username_, email_ from users_`

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_"}).
		AddRow(uint(23), false, 69)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	users, err := repo.GetAll()

	require.Error(t, err, "should return scan error")

	require.Nil(t, users, "should not return users")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoGetByUsername(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select id_, username_, email_ from users_ where username_ = $1`

	mockRows := mock.
		NewRows([]string{"id_", "username_", "email_"}).
		AddRow(uint(23), "janedoe", "jane@example.org")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("janedoe").WillReturnRows(mockRows)

	user, err := repo.GetByAttribute("username", "janedoe")

	require.NoError(t, err, "should not return error")

	require.Equal(t, &User{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, user, "should return matching user")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoGetByUnknownAttribute(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	user, err := repo.GetByAttribute("foo", "bar")

	require.Error(t, err, "should return error")
	require.Nil(t, user, "should not return user")
}

func testUserRepoGetByAttributeDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select id_, username_, email_ from users_ where username_ = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("janedoe").WillReturnError(errors.New("db_error"))

	user, err := repo.GetByAttribute("username", "janedoe")

	require.Error(t, err, "should return error")
	require.Nil(t, user, "should not return user")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoGetPasswordByUsername(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select password_ from users_ where username_ = $1`

	mockRow := mock.NewRows([]string{"password_"}).AddRow("p4ssw0rd")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("janedoe").WillReturnRows(mockRow)

	password, err := repo.GetPasswordByUsername("janedoe")

	require.NoError(t, err, "should not return error")
	require.Equal(t, "p4ssw0rd", password, "should return hashed password")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()
}

func testUserRepoGetPasswordByUsernameDbError(t *testing.T, mock pgxmock.PgxPoolIface, repo userPostgresRepository) {
	query := `select password_ from users_ where username_ = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("janedoe").WillReturnError(errors.New("db_error"))

	password, err := repo.GetPasswordByUsername("janedoe")

	require.EqualError(t, err, "db_error", "should not return error")
	require.Empty(t, password, "should not return password")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("expectations not met: ", err)
	}

	mock.Reset()

}
