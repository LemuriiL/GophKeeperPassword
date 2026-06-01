package model

import "time"

const (
	TypeLoginPassword = "login_password"
	TypeText          = "text"
	TypeBinary        = "binary"
	TypeCard          = "card"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
}

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
