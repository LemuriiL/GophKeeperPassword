package dto

import "time"

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

type UpsertItemRequest struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Meta       string `json:"meta"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Salt       string `json:"salt"`
}

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
