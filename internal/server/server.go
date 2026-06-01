package server

import (
	"context"
	"net/http"
	"time"
)

// App хранит зависимости сервера.
type App struct {
	cfg     Config
	store   *SQLite
	http    *http.Server
	handler *Handler
}

// NewApp создаёт серверное приложение.
func NewApp(cfg Config) (*App, error) {
	store, err := NewSQLite(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	auth := NewAuthService(store)
	items := NewItemService(store)
	tokens := NewTokenManager(cfg.JWTSecret)
	h := NewHandler(auth, items, tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.Handle("POST /api/items", AuthMiddleware(tokens, http.HandlerFunc(h.UpsertItem)))
	mux.Handle("GET /api/items", AuthMiddleware(tokens, http.HandlerFunc(h.ListItems)))
	mux.Handle("DELETE /api/items/", AuthMiddleware(tokens, http.HandlerFunc(h.DeleteItem)))

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           LoggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &App{cfg: cfg, store: store, http: srv, handler: h}, nil
}

// Run запускает сервер.
func (a *App) Run() error {
	err := a.http.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown останавливает сервер.
func (a *App) Shutdown(ctx context.Context) error {
	return a.http.Shutdown(ctx)
}

// Close закрывает базу.
func (a *App) Close() error {
	return a.store.Close()
}
