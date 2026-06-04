package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("CONFIG", "")
	t.Setenv("ADDRESS", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("JWT_SECRET", "")

	cfg, err := LoadConfig("", "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Address != ":8080" {
		t.Fatalf("unexpected address: %s", cfg.Address)
	}

	if cfg.DBPath != "gophkeeper.db" {
		t.Fatalf("unexpected db path: %s", cfg.DBPath)
	}

	if cfg.JWTSecret != "supersecret" {
		t.Fatalf("unexpected secret: %s", cfg.JWTSecret)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:   ":9090",
		DBPath:    "test.db",
		JWTSecret: "file-secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path, data, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("CONFIG", "")
	t.Setenv("ADDRESS", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("JWT_SECRET", "")

	cfg, err := LoadConfig(path, "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Address != ":9090" {
		t.Fatalf("unexpected address: %s", cfg.Address)
	}

	if cfg.DBPath != "test.db" {
		t.Fatalf("unexpected db path: %s", cfg.DBPath)
	}

	if cfg.JWTSecret != "file-secret" {
		t.Fatalf("unexpected secret: %s", cfg.JWTSecret)
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:   ":9090",
		DBPath:    "test.db",
		JWTSecret: "file-secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path, data, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("ADDRESS", ":7070")
	t.Setenv("DB_PATH", "env.db")
	t.Setenv("JWT_SECRET", "env-secret")
	t.Setenv("CONFIG", "")

	cfg, err := LoadConfig(path, "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Address != ":7070" {
		t.Fatalf("unexpected address: %s", cfg.Address)
	}

	if cfg.DBPath != "env.db" {
		t.Fatalf("unexpected db path: %s", cfg.DBPath)
	}

	if cfg.JWTSecret != "env-secret" {
		t.Fatalf("unexpected secret: %s", cfg.JWTSecret)
	}
}

func TestLoadConfigBadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")

	err := os.WriteFile(path, []byte(`{"address":`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(path, "")
	if err == nil {
		t.Fatal("expected error")
	}
}
