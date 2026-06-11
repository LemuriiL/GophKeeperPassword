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

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		out.Address = v
	}

	if v, ok := os.LookupEnv("DB_PATH"); ok {
		out.DBPath = v
	}

	if v, ok := os.LookupEnv("JWT_SECRET"); ok {
		out.JWTSecret = v
	}

	if v, ok := os.LookupEnv("TLS_CERT_FILE"); ok {
		out.TLSCertFile = v
	}

	if v, ok := os.LookupEnv("TLS_KEY_FILE"); ok {
		out.TLSKeyFile = v
	}

	if strings.TrimSpace(out.JWTSecret) == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	if strings.TrimSpace(out.TLSCertFile) == "" {
		return Config{}, errors.New("TLS_CERT_FILE is required")
	}

	if strings.TrimSpace(out.TLSKeyFile) == "" {
		return Config{}, errors.New("TLS_KEY_FILE is required")
	}

	return out, nil
}

// pickString выбирает первое непустое строковое значение
func pickString(envName string, values ...string) string {
	if v, ok := os.LookupEnv(envName); ok {
		return v
	}

	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
