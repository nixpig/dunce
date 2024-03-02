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

type Claims struct {
	UserId   int    `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *User) (string, error) {
	claims := Claims{
		UserId:   user.Id,
		UserRole: user.Role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(config.Get("SECRET")))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(config.Get("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
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
		return "", fmt.Errorf("user is already logged out and unable to logout")
	}

	signedToken, err := GenerateToken(&User{Id: userId, UserData: UserData{Role: userRole}})
	if err != nil {
		return "", err
	}

	session := &Session{
		Token:  signedToken,
		UserId: userId,
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
