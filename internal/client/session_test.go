package client

import (
	"path/filepath"
	"testing"
)

// TestSaveLoadSession проверяет сохранение и чтение сессии
func TestSaveLoadSession(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")

	err := SaveSession(path, Session{
		Login: "user1",
		Token: "token1",
	})
	if err != nil {
		t.Fatal(err)
	}

	session, err := LoadSession(path)
	if err != nil {
		t.Fatal(err)
	}

	if session.Login != "user1" {
		t.Fatalf("unexpected login: %s", session.Login)
	}

	if session.Token != "token1" {
		t.Fatalf("unexpected token: %s", session.Token)
	}
}

// TestLoadSessionMissingFile проверяет ошибку при отсутствии файла сессии
func TestLoadSessionMissingFile(t *testing.T) {
	_, err := LoadSession(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

// TestSaveSessionBadPath проверяет ошибку записи сессии
func TestSaveSessionBadPath(t *testing.T) {
	err := SaveSession(filepath.Join(t.TempDir(), "missing", "session.json"), Session{
		Login: "user1",
		Token: "token1",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
