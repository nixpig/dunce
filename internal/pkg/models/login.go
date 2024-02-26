package models

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nixpig/dunce/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Db Dbconn
}

type LoginDetails struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=255"`
}

type Session struct {
	Token     string
	UserId    int
	IssuedAt  int64
	ExpiresAt int64
}

func (l *Login) WithUsernamePassword(user *LoginDetails) (string, error) {
	getUserQuery := `select password_ from users_ where username_ = $1`
	row := l.Db.QueryRow(context.Background(), getUserQuery, user.Username)

	var savedHash string

	if err := row.Scan(&savedHash); err != nil {
		return "", fmt.Errorf("invalid login details")
	}

	if !comparePasswordHash(savedHash, user.Password) {
		return "", fmt.Errorf("invalid login details")
	}

	userDataQuery := `select id_, role_ from users_ where username_ = $1`
	var userId int
	var userRole RoleName

	userDataRow := l.Db.QueryRow(context.Background(), userDataQuery, user.Username)

	if err := userDataRow.Scan(&userId, &userRole); err != nil {
		return "", fmt.Errorf("unable to login")
	}

	if err := l.Logout(userId); err != nil {
		return "", fmt.Errorf("user is already logged in and unable to logout")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	signedToken, err := token.SignedString([]byte(config.Get("SECRET")))
	if err != nil {
		return "", fmt.Errorf("unable to login")
	}

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(time.Hour * 24).Unix()

	claims["user_id"] = userId
	claims["user_role"] = userRole
	claims["iat"] = issuedAt
	claims["exp"] = expiresAt

	session := &Session{
		Token:     signedToken,
		UserId:    userId,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	// save this new session for user
	createSessionQuery := `insert into sessions_ (token_, user_id_, issued_at_, expires_at_) values ($1, $2, $3, $4)`
	createSessionRes, err := l.Db.Exec(context.Background(), createSessionQuery, &session.Token, &session.UserId, &session.IssuedAt, &session.ExpiresAt)
	if err != nil {
		return "", fmt.Errorf("unable to create session")
	}

	if createSessionRes.RowsAffected() == 0 {
		return "", fmt.Errorf("unable to save session")
	}

	return signedToken, nil
}

func (l *Login) Logout(userId int) error {
	logoutQuery := `delete from sessions_ where user_id_ = $1`
	_, _ = l.Db.Exec(context.Background(), logoutQuery, userId)

	return nil
}

func comparePasswordHash(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
