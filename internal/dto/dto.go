package dto

import "time"

// RegisterRequest описывает регистрацию
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest описывает логин
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse возвращает токен
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginPasswordSecret хранит пару логин пароль
type LoginPasswordSecret struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// TextSecret хранит текстовый секрет
type TextSecret struct {
	Text string `json:"text"`
}

// BinarySecret хранит бинарные данные
type BinarySecret struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// CardSecret хранит данные карты
type CardSecret struct {
	Number   string `json:"number"`
	Holder   string `json:"holder"`
	Expiry   string `json:"expiry"`
	CVV      string `json:"cvv"`
	Provider string `json:"provider"`
}

// UpsertItemRequest описывает создание или обновление секрета
type UpsertItemRequest struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Meta       string `json:"meta"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Salt       string `json:"salt"`
}

// ItemResponse описывает ответ с секретом
type ItemResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Meta       string    `json:"meta"`
	Ciphertext string    `json:"ciphertext"`
	Nonce      string    `json:"nonce"`
	Salt       string    `json:"salt"`
	UpdatedAt  time.Time `json:"updated_at"`
}
