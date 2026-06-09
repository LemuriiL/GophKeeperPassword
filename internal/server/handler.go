package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

var ErrLoginExists = errors.New("login already exists")

// Handler хранит HTTP хендлеры сервера
type Handler struct {
	auth   *AuthService
	items  *ItemService
	tokens *TokenManager
}

// NewHandler создает набор хендлеров
func NewHandler(auth *AuthService, items *ItemService, tokens *TokenManager) *Handler {
	return &Handler{
		auth:   auth,
		items:  items,
		tokens: tokens,
	}
}

// Register регистрирует пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	login := strings.TrimSpace(req.Login)
	if login == "" || req.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return
	}

	user, err := h.auth.Register(r.Context(), login, req.Password)
	if err != nil {
		if errors.Is(err, ErrLoginExists) {
			http.Error(w, "login already exists", http.StatusConflict)
			return
		}

		slog.Error("register user", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := h.tokens.Issue(user.ID, user.Login)
	if err != nil {
		slog.Error("issue token after register", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dto.LoginResponse{Token: token})
}

// Login логинит пользователя
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	login := strings.TrimSpace(req.Login)
	if login == "" || req.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return
	}

	user, err := h.auth.Login(r.Context(), login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.tokens.Issue(user.ID, user.Login)
	if err != nil {
		slog.Error("issue token after login", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.LoginResponse{Token: token})
}

// UpsertItem создает или обновляет секрет
func (h *Handler) UpsertItem(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Type) == "" || strings.TrimSpace(req.Title) == "" {
		http.Error(w, "empty type or title", http.StatusBadRequest)
		return
	}

	id := req.ID
	if r.Method == http.MethodPut {
		id = strings.TrimPrefix(r.URL.Path, "/api/items/")
		if strings.TrimSpace(id) == "" {
			http.Error(w, "empty id", http.StatusBadRequest)
			return
		}
	}

	item, err := h.items.Save(r.Context(), userIDFromContext(r.Context()), model.Item{
		ID:         id,
		Type:       req.Type,
		Title:      req.Title,
		Meta:       req.Meta,
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
		Salt:       req.Salt,
	})
	if err != nil {
		slog.Error("upsert item", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toItemResponse(item))
}

// GetItem возвращает один секрет
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if strings.TrimSpace(id) == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	item, err := h.items.Get(r.Context(), userIDFromContext(r.Context()), id)
	if err != nil {
		if IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		slog.Error("get item", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toItemResponse(item))
}

// ListItems возвращает список секретов
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.items.List(r.Context(), userIDFromContext(r.Context()))
	if err != nil {
		slog.Error("list items", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.ItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toItemResponse(item))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// DeleteItem удаляет секрет
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if strings.TrimSpace(id) == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	if err := h.items.Delete(r.Context(), userIDFromContext(r.Context()), id); err != nil {
		if IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		slog.Error("delete item", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toItemResponse(item model.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:         item.ID,
		Type:       item.Type,
		Title:      item.Title,
		Meta:       item.Meta,
		Ciphertext: item.Ciphertext,
		Nonce:      item.Nonce,
		Salt:       item.Salt,
		UpdatedAt:  item.UpdatedAt,
	}
}
