package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func userIDFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(userIDKey).(int64)
	return v
}

func withUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// AuthMiddleware проверяет Bearer токен
func AuthMiddleware(tokens *TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		if !strings.HasPrefix(raw, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(raw, "Bearer ")
		userID, _, err := tokens.Parse(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(withUserID(r.Context(), userID)))
	})
}

// LoggingMiddleware логирует HTTP запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{
			ResponseWriter: w,
			code:           http.StatusOK,
		}

		next.ServeHTTP(sw, r)

		slog.Info(
			"request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.code,
			"duration", time.Since(start),
		)
	})
}
