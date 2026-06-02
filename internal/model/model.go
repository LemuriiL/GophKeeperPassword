package model

import "time"

// TypeLoginPassword это тип секрета логин пароль
const TypeLoginPassword = "login_password"

// TypeText это тип текстового секрета
const TypeText = "text"

// TypeBinary это тип бинарного секрета
const TypeBinary = "binary"

// TypeCard это тип данных карты
const TypeCard = "card"

// User описывает пользователя
type User struct {
	ID           int64
	Login        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
}

// Item описывает секрет пользователя
type Item struct {
	ID         string
	UserID     int64
	Type       string
	Title      string
	Meta       string
	Ciphertext string
	Nonce      string
	Salt       string
	UpdatedAt  time.Time
}
