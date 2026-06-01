package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

type Handler struct {
	auth   *AuthService
	items  *ItemService
	tokens *TokenManager
}

func NewHandler(auth *AuthService, items *ItemService, tokens *TokenManager) *Handler {
	return &Handler{
		auth:   auth,
		items:  items,
		tokens: tokens,
	}
}

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

	ok, err := h.auth.IsUniqueLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "login already exists", http.StatusConflict)
		return
	}

	if _, err := h.auth.Register(r.Context(), login, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.LoginResponse{
		Token: token,
	})
}

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

	item, err := h.items.Save(r.Context(), userIDFromContext(r.Context()), model.Item{
		Type:       req.Type,
		Title:      req.Title,
		Meta:       req.Meta,
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ItemResponse{
		ID:         item.ID,
		Type:       item.Type,
		Title:      item.Title,
		Meta:       item.Meta,
		Ciphertext: item.Ciphertext,
		Nonce:      item.Nonce,
		UpdatedAt:  item.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.items.List(r.Context(), userIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.ItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.ItemResponse{
			ID:         item.ID,
			Type:       item.Type,
			Title:      item.Title,
			Meta:       item.Meta,
			Ciphertext: item.Ciphertext,
			Nonce:      item.Nonce,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if strings.TrimSpace(id) == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	if err := h.items.Delete(r.Context(), userIDFromContext(r.Context()), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
