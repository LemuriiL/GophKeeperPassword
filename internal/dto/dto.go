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
}

type LoginPasswordSecret struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TextSecret struct {
	Text string `json:"text"`
}

type BinarySecret struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type CardSecret struct {
	Number   string `json:"number"`
	Holder   string `json:"holder"`
	Expiry   string `json:"expiry"`
	CVV      string `json:"cvv"`
	Provider string `json:"provider"`
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
