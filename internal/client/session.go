package client

import (
	"encoding/json"
	"os"
)

// Session хранит локальную сессию
type Session struct {
	Login string `json:"login"`
	Token string `json:"token"`
}

// SaveSession сохраняет сессию на диск
func SaveSession(path string, session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

// LoadSession читает сессию с диска
func LoadSession(path string) (Session, error) {
	var session Session

	data, err := os.ReadFile(path)
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(data, &session)
	return session, err
}
