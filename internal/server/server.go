package server

import (
	"context"
	"net/http"
	"time"
)

type App struct {
	cfg     Config
	store   *SQLite
	http    *http.Server
	handler *Handler
	tokens  *TokenManager
}

func NewApp(cfg Config) (*App, error) {
	store, err := NewSQLite(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	auth := NewAuthService(store)
	items := NewItemService(store)
	tokens := NewTokenManager(cfg.JWTSecret)
	h := NewHandler(auth, items, tokens)

	app := &App{
		cfg:     cfg,
		store:   store,
		handler: h,
		tokens:  tokens,
	}

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           LoggingMiddleware(app.routes()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	app.http = srv

	return app, nil
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", a.handler.Register)
	mux.HandleFunc("POST /api/login", a.handler.Login)
	mux.Handle("POST /api/items", AuthMiddleware(a.tokens, http.HandlerFunc(a.handler.UpsertItem)))
	mux.Handle("GET /api/items", AuthMiddleware(a.tokens, http.HandlerFunc(a.handler.ListItems)))
	mux.Handle("GET /api/items/", AuthMiddleware(a.tokens, http.HandlerFunc(a.handler.GetItem)))
	mux.Handle("PUT /api/items/", AuthMiddleware(a.tokens, http.HandlerFunc(a.handler.UpsertItem)))
	mux.Handle("DELETE /api/items/", AuthMiddleware(a.tokens, http.HandlerFunc(a.handler.DeleteItem)))

	return mux
}

func (a *App) RoutesForTests() http.Handler {
	return a.routes()
}

func (a *App) Run() error {
	err := a.http.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.http.Shutdown(ctx)
}

func (a *App) Close() error {
	return a.store.Close()
}
