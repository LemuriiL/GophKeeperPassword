package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// AuthMiddleware проверяет JWT.
func AuthMiddleware(tokens *TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, _, err := tokens.Parse(parts[1])
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware пишет лог запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}

func userIDFromContext(ctx context.Context) int64 {
	id, _ := ctx.Value(userIDKey).(int64)
	return id
}
