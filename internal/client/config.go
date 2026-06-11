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
	ServerURL          string `json:"server_url"`
	SessionFile        string `json:"session_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

// LoadConfig загружает конфиг клиента
func LoadConfig(shortPath string, longPath string) (Config, error) {
	path := pickConfigPath(shortPath, longPath)

	home, _ := os.UserHomeDir()
	defaultSession := filepath.Join(home, ".gophkeeper_session.json")

	out := Config{
		ServerURL:   "https://localhost:8080",
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

		out.InsecureSkipVerify = fileCfg.InsecureSkipVerify
	}

	if v, ok := os.LookupEnv("SERVER_URL"); ok {
		out.ServerURL = v
	}

	if v, ok := os.LookupEnv("SESSION_FILE"); ok {
		out.SessionFile = v
	}

	if v, ok := os.LookupEnv("INSECURE_SKIP_VERIFY"); ok {
		out.InsecureSkipVerify = v == "true" || v == "1"
	}

	if runtime.GOOS == "windows" && strings.HasPrefix(out.SessionFile, "~") {
		out.SessionFile = strings.Replace(out.SessionFile, "~", home, 1)
	}

	return out, nil
}

// pickConfigPath выбирает путь к конфигу клиента
func pickConfigPath(shortPath string, longPath string) string {
	if v, ok := os.LookupEnv("CONFIG"); ok {
		return v
	}

	if strings.TrimSpace(shortPath) != "" {
		return shortPath
	}

	return longPath
}
