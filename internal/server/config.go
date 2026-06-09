package server

import (
	"errors"
	"os"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/cfg"
)

// Config хранит настройки сервера
type Config struct {
	Address     string `json:"address"`
	DBPath      string `json:"db_path"`
	JWTSecret   string `json:"jwt_secret"`
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`
}

// LoadConfig загружает конфиг сервера
func LoadConfig(shortPath string, longPath string) (Config, error) {
	path := pickString("CONFIG", shortPath, longPath)

	out := Config{
		Address: ":8080",
		DBPath:  "gophkeeper.db",
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

		if fileCfg.TLSCertFile != "" {
			out.TLSCertFile = fileCfg.TLSCertFile
		}

		if fileCfg.TLSKeyFile != "" {
			out.TLSKeyFile = fileCfg.TLSKeyFile
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

	if v := os.Getenv("TLS_CERT_FILE"); v != "" {
		out.TLSCertFile = v
	}

	if v := os.Getenv("TLS_KEY_FILE"); v != "" {
		out.TLSKeyFile = v
	}

	if strings.TrimSpace(out.JWTSecret) == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	if strings.TrimSpace(out.TLSCertFile) == "" && strings.TrimSpace(out.TLSKeyFile) != "" {
		return Config{}, errors.New("TLS_CERT_FILE is required when TLS_KEY_FILE is set")
	}

	if strings.TrimSpace(out.TLSCertFile) != "" && strings.TrimSpace(out.TLSKeyFile) == "" {
		return Config{}, errors.New("TLS_KEY_FILE is required when TLS_CERT_FILE is set")
	}

	return out, nil
}

// pickString выбирает первое непустое строковое значение
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
