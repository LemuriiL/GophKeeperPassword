package dto

import "time"

// RegisterRequest описывает регистрацию пользователя.
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest описывает логин пользователя.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse возвращает токен.
type LoginResponse struct {
	Token string `json:"token"`
}

// UpsertItemRequest описывает сохранение секрета.
type UpsertItemRequest struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Meta       string `json:"meta"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

// ItemResponse описывает ответ с секретом.
type ItemResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Meta       string    `json:"meta"`
	Ciphertext string    `json:"ciphertext"`
	Nonce      string    `json:"nonce"`
	UpdatedAt  time.Time `json:"updated_at"`
}
