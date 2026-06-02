package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "client.json")
	sessionPath := filepath.Join(dir, "session.json")

	cfgData, err := json.Marshal(map[string]string{
		"server_url":   "http://example.com",
		"session_file": sessionPath,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(cfgPath, cfgData, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(cfgPath, "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "http://example.com" {
		t.Fatalf("unexpected server url: %s", cfg.ServerURL)
	}

	if cfg.SessionFile != sessionPath {
		t.Fatalf("unexpected session file: %s", cfg.SessionFile)
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	t.Setenv("SERVER_URL", "http://env.example.com")
	t.Setenv("SESSION_FILE", "env-session.json")

	cfg, err := LoadConfig("", "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "http://env.example.com" {
		t.Fatalf("unexpected server url: %s", cfg.ServerURL)
	}

	if cfg.SessionFile != "env-session.json" {
		t.Fatalf("unexpected session file: %s", cfg.SessionFile)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("SERVER_URL", "")
	t.Setenv("SESSION_FILE", "")
	t.Setenv("CONFIG", "")

	cfg, err := LoadConfig("", "")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(cfg.ServerURL, "http://") {
		t.Fatalf("unexpected default server url: %s", cfg.ServerURL)
	}

	if strings.TrimSpace(cfg.SessionFile) == "" {
		t.Fatal("expected default session file")
	}
}
