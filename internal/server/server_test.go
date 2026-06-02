package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp(Config{
		Address:   ":8080",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if app.store == nil {
		t.Fatal("expected store")
	}

	if app.handler == nil {
		t.Fatal("expected handler")
	}

	if app.tokens == nil {
		t.Fatal("expected tokens")
	}
}

func TestRoutesUnauthorized(t *testing.T) {
	app, err := NewApp(Config{
		Address:   ":8080",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	rec := httptest.NewRecorder()

	app.RoutesForTests().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestRoutesNotFound(t *testing.T) {
	app, err := NewApp(Config{
		Address:   ":8080",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	req := httptest.NewRequest(http.MethodGet, "/no-such-route", nil)
	rec := httptest.NewRecorder()

	app.RoutesForTests().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestShutdownWithoutRun(t *testing.T) {
	app, err := NewApp(Config{
		Address:   ":8080",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if err = app.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}
