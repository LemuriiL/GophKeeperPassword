package server

import (
	"context"
	"testing"
)

func TestAuthServiceRegisterLogin(t *testing.T) {
	store, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	svc := NewAuthService(store)

	user, err := svc.Register(context.Background(), "user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatal("expected user id")
	}

	loggedIn, err := svc.Login(context.Background(), "user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if loggedIn.Login != "user1" {
		t.Fatalf("unexpected login: %s", loggedIn.Login)
	}
}

func TestAuthServiceIsUniqueLogin(t *testing.T) {
	store, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	svc := NewAuthService(store)

	ok, err := svc.IsUniqueLogin(context.Background(), "user1")
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected unique login")
	}

	_, err = svc.Register(context.Background(), "user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	ok, err = svc.IsUniqueLogin(context.Background(), "user1")
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("expected non-unique login")
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	store, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	svc := NewAuthService(store)

	_, err = svc.Register(context.Background(), "user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Login(context.Background(), "user1", "wrong")
	if err == nil {
		t.Fatal("expected login error")
	}
}
