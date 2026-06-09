package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/model"
	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
)

// AuthService отвечает за регистрацию и логин
type AuthService struct {
	store *SQLite
}

// NewAuthService создает сервис авторизации
func NewAuthService(store *SQLite) *AuthService {
	return &AuthService{store: store}
}

// Register создает нового пользователя
func (s *AuthService) Register(ctx context.Context, login string, password string) (model.User, error) {
	saltRaw := make([]byte, 16)
	if _, err := rand.Read(saltRaw); err != nil {
		return model.User{}, err
	}

	hash := secure.DeriveKey(password, saltRaw)

	user := model.User{
		Login:        login,
		PasswordHash: base64.StdEncoding.EncodeToString(hash),
		Salt:         base64.StdEncoding.EncodeToString(saltRaw),
		CreatedAt:    Now(),
	}

	id, err := s.store.CreateUser(ctx, user)
	if err != nil {
		if isUniqueLoginErr(err) {
			return model.User{}, ErrLoginExists
		}
		return model.User{}, err
	}

	user.ID = id
	return user, nil
}

// Login проверяет логин и пароль
func (s *AuthService) Login(ctx context.Context, login string, password string) (model.User, error) {
	user, err := s.store.GetUserByLogin(ctx, login)
	if err != nil {
		return model.User{}, err
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return model.User{}, err
	}

	hash := secure.DeriveKey(password, salt)
	if base64.StdEncoding.EncodeToString(hash) != user.PasswordHash {
		return model.User{}, errors.New("invalid login or password")
	}

	return user, nil
}

// IsUniqueLogin проверяет свободен ли логин
func (s *AuthService) IsUniqueLogin(ctx context.Context, login string) (bool, error) {
	_, err := s.store.GetUserByLogin(ctx, login)
	if err == nil {
		return false, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}

	return false, err
}

func isUniqueLoginErr(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint failed")
}
