package model

import "time"

// User описывает пользователя.
type User struct {
	ID           int64
	Login        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
}

// Item описывает секрет пользователя.
type Item struct {
	ID         string
	UserID     int64
	Type       string
	Title      string
	Meta       string
	Ciphertext string
	Nonce      string
	UpdatedAt  time.Time
}
