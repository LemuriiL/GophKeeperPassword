package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	err := os.WriteFile(path, []byte(`{"name":"demo","count":7}`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load[testConfig](path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "demo" {
		t.Fatalf("unexpected name: %s", cfg.Name)
	}

	if cfg.Count != 7 {
		t.Fatalf("unexpected count: %d", cfg.Count)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load[testConfig](filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadBadJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")

	err := os.WriteFile(path, []byte(`{"name":`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load[testConfig](path)
	if err == nil {
		t.Fatal("expected error")
	}
}
