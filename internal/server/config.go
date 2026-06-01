package server

import (
	"os"
	"strconv"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/cfg"
)

// Config хранит настройки сервера.
type Config struct {
	Address   string `json:"address"`
	DBPath    string `json:"db_path"`
	JWTSecret string `json:"jwt_secret"`
}

// LoadConfig загружает конфиг сервера.
func LoadConfig(shortPath string, longPath string) (Config, error) {
	path := pickString("CONFIG", shortPath, longPath)
	out := Config{
		Address:   ":8080",
		DBPath:    "gophkeeper.db",
		JWTSecret: "supersecret",
	}

	if strings.TrimSpace(path) != "" {
		fileCfg, err := cfg.Load[Config](path)
		if err != nil {
			return Config{}, err
		}
		if fileCfg.Address != "" {
			out.Address = fileCfg.Address
		}
		if fileCfg.DBPath != "" {
			out.DBPath = fileCfg.DBPath
		}
		if fileCfg.JWTSecret != "" {
			out.JWTSecret = fileCfg.JWTSecret
		}
	}

	if v := os.Getenv("ADDRESS"); v != "" {
		out.Address = v
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		out.DBPath = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		out.JWTSecret = v
	}

	return out, nil
}

func pickString(envName string, values ...string) string {
	if v := os.Getenv(envName); v != "" {
		return v
	}

	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func pickInt(envName string, def int) int {
	if v := os.Getenv(envName); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}

	return def
}
