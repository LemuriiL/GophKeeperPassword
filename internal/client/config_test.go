package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func unsetClientConfigEnv(t *testing.T) {
	t.Helper()

	keys := []string{
		"CONFIG",
		"SERVER_URL",
		"SESSION_FILE",
		"INSECURE_SKIP_VERIFY",
	}

	oldValues := make(map[string]string, len(keys))
	oldExists := make(map[string]bool, len(keys))

	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		oldValues[key] = value
		oldExists[key] = ok
		_ = os.Unsetenv(key)
	}

	t.Cleanup(func() {
		for _, key := range keys {
			if oldExists[key] {
				_ = os.Setenv(key, oldValues[key])
			} else {
				_ = os.Unsetenv(key)
			}
		}
	})
}

func TestLoadConfigFromFile(t *testing.T) {
	unsetClientConfigEnv(t)

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
	unsetClientConfigEnv(t)

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
	unsetClientConfigEnv(t)

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

func TestLoadConfigEmptyEnvOverride(t *testing.T) {
	unsetClientConfigEnv(t)

	t.Setenv("SERVER_URL", "")

	cfg, err := LoadConfig("", "")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.ServerURL != "" {
		t.Fatalf("unexpected server url: %s", cfg.ServerURL)
	}
}
