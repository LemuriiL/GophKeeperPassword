package server

import "testing"

func TestTokenManagerIssueParse(t *testing.T) {
	m := NewTokenManager("secret")

	token, err := m.Issue(42, "user")
	if err != nil {
		t.Fatal(err)
	}

	userID, login, err := m.Parse(token)
	if err != nil {
		t.Fatal(err)
	}

	if userID != 42 {
		t.Fatalf("unexpected user id: %d", userID)
	}

	if login != "user" {
		t.Fatalf("unexpected login: %s", login)
	}
}

func TestTokenManagerParseInvalid(t *testing.T) {
	m := NewTokenManager("secret")

	_, _, err := m.Parse("bad-token")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
