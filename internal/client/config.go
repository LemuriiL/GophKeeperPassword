package client

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/cfg"
)

// Config хранит настройки клиента
type Config struct {
	ServerURL   string `json:"server_url"`
	SessionFile string `json:"session_file"`
}

// LoadConfig загружает конфиг клиента
func LoadConfig(shortPath string, longPath string) (Config, error) {
	path := os.Getenv("CONFIG")
	if path == "" {
		if strings.TrimSpace(shortPath) != "" {
			path = shortPath
		} else {
			path = longPath
		}
	}

	home, _ := os.UserHomeDir()
	defaultSession := filepath.Join(home, ".gophkeeper_session.json")

	out := Config{
		ServerURL:   "http://localhost:8080",
		SessionFile: defaultSession,
	}

	if strings.TrimSpace(path) != "" {
		fileCfg, err := cfg.Load[Config](path)
		if err != nil {
			return Config{}, err
		}

		if fileCfg.ServerURL != "" {
			out.ServerURL = fileCfg.ServerURL
		}

		if fileCfg.SessionFile != "" {
			out.SessionFile = fileCfg.SessionFile
		}
	}

	if v := os.Getenv("SERVER_URL"); v != "" {
		out.ServerURL = v
	}

	if v := os.Getenv("SESSION_FILE"); v != "" {
		out.SessionFile = v
	}

	if runtime.GOOS == "windows" && strings.HasPrefix(out.SessionFile, "~") {
		out.SessionFile = strings.Replace(out.SessionFile, "~", home, 1)
	}

	return out, nil
}
