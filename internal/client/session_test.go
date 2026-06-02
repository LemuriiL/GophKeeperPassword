package client

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadSession(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")

	err := SaveSession(path, Session{
		Login: "user1",
		Token: "token1",
		Salt:  "salt1",
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

	if session.Salt != "salt1" {
		t.Fatalf("unexpected salt: %s", session.Salt)
	}
}

func TestLoadSessionMissingFile(t *testing.T) {
	_, err := LoadSession(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}
