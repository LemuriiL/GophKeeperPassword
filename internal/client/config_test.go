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

	cfgData, err := json.Marshal(map[string]any{
		"server_url":           "http://example.com",
		"session_file":         sessionPath,
		"insecure_skip_verify": true,
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

	if !cfg.InsecureSkipVerify {
		t.Fatal("expected insecure skip verify")
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	t.Setenv("SERVER_URL", "http://env.example.com")
	t.Setenv("SESSION_FILE", "env-session.json")
	t.Setenv("INSECURE_SKIP_VERIFY", "true")

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

	if !cfg.InsecureSkipVerify {
		t.Fatal("expected insecure skip verify")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("SERVER_URL", "")
	t.Setenv("SESSION_FILE", "")
	t.Setenv("CONFIG", "")
	t.Setenv("INSECURE_SKIP_VERIFY", "")

	cfg, err := LoadConfig("", "")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(cfg.ServerURL, "https://") {
		t.Fatalf("unexpected default server url: %s", cfg.ServerURL)
	}

	if strings.TrimSpace(cfg.SessionFile) == "" {
		t.Fatal("expected default session file")
	}

	if cfg.InsecureSkipVerify {
		t.Fatal("unexpected insecure skip verify")
	}
}
