package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaultsRequiresJWTSecret(t *testing.T) {
	t.Setenv("CONFIG", "")
	t.Setenv("ADDRESS", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("TLS_CERT_FILE", "")
	t.Setenv("TLS_KEY_FILE", "")

	_, err := LoadConfig("", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:     ":9090",
		DBPath:      "test.db",
		JWTSecret:   "file-secret",
		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",
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
	t.Setenv("TLS_CERT_FILE", "")
	t.Setenv("TLS_KEY_FILE", "")

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

	if cfg.TLSCertFile != "cert.pem" {
		t.Fatalf("unexpected tls cert file: %s", cfg.TLSCertFile)
	}

	if cfg.TLSKeyFile != "key.pem" {
		t.Fatalf("unexpected tls key file: %s", cfg.TLSKeyFile)
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:     ":9090",
		DBPath:      "test.db",
		JWTSecret:   "file-secret",
		TLSCertFile: "file-cert.pem",
		TLSKeyFile:  "file-key.pem",
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
	t.Setenv("TLS_CERT_FILE", "env-cert.pem")
	t.Setenv("TLS_KEY_FILE", "env-key.pem")
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

	if cfg.TLSCertFile != "env-cert.pem" {
		t.Fatalf("unexpected tls cert file: %s", cfg.TLSCertFile)
	}

	if cfg.TLSKeyFile != "env-key.pem" {
		t.Fatalf("unexpected tls key file: %s", cfg.TLSKeyFile)
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

func TestLoadConfigTLSCertWithoutKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:     ":9090",
		DBPath:      "test.db",
		JWTSecret:   "file-secret",
		TLSCertFile: "cert.pem",
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
	t.Setenv("TLS_CERT_FILE", "")
	t.Setenv("TLS_KEY_FILE", "")

	_, err = LoadConfig(path, "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfigTLSKeyWithoutCert(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.json")

	data, err := json.Marshal(Config{
		Address:    ":9090",
		DBPath:     "test.db",
		JWTSecret:  "file-secret",
		TLSKeyFile: "key.pem",
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
	t.Setenv("TLS_CERT_FILE", "")
	t.Setenv("TLS_KEY_FILE", "")

	_, err = LoadConfig(path, "")
	if err == nil {
		t.Fatal("expected error")
	}
}
