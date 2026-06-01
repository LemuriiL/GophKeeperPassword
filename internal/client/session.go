package client

import (
	"encoding/json"
	"os"
)

type Session struct {
	Login string `json:"login"`
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

func SaveSession(path string, session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func LoadSession(path string) (Session, error) {
	var session Session

	data, err := os.ReadFile(path)
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(data, &session)
	return session, err
}
