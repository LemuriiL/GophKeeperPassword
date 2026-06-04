package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager работает с JWT токенами
type TokenManager struct {
	secret []byte
}

// NewTokenManager создает менеджер токенов
func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: []byte(secret)}
}

// Issue выпускает токен для пользователя
func (m *TokenManager) Issue(userID int64, login string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"login": login,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse разбирает токен и возвращает данные пользователя
func (m *TokenManager) Parse(raw string) (int64, string, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	idf, _ := claims["sub"].(float64)
	login, _ := claims["login"].(string)

	return int64(idf), login, nil
}
