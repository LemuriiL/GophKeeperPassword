package client

import (
	"encoding/base64"
	"errors"

	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
)

// EncryptPayload шифрует строку мастер-ключом.
func EncryptPayload(password string, salt string, plaintext string) (string, string, error) {
	saltRaw, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", "", err
	}

	key := secure.DeriveKey(password, saltRaw)
	return secure.Encrypt(key, []byte(plaintext))
}

// DecryptPayload расшифровывает строку мастер-ключом.
func DecryptPayload(password string, salt string, ciphertext string, nonce string) (string, error) {
	saltRaw, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	key := secure.DeriveKey(password, saltRaw)
	plain, err := secure.Decrypt(key, ciphertext, nonce)
	if err != nil {
		return "", err
	}
	if plain == nil {
		return "", errors.New("empty plaintext")
	}

	return string(plain), nil
}
