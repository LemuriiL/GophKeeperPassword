package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

func newTestHTTP(t *testing.T) (*Handler, *TokenManager, http.Handler, func()) {
	t.Helper()

	store, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	auth := NewAuthService(store)
	items := NewItemService(store)
	tokens := NewTokenManager("secret")
	h := NewHandler(auth, items, tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.Handle("POST /api/items", AuthMiddleware(tokens, http.HandlerFunc(h.UpsertItem)))
	mux.Handle("GET /api/items", AuthMiddleware(tokens, http.HandlerFunc(h.ListItems)))
	mux.Handle("GET /api/items/", AuthMiddleware(tokens, http.HandlerFunc(h.GetItem)))
	mux.Handle("PUT /api/items/", AuthMiddleware(tokens, http.HandlerFunc(h.UpsertItem)))
	mux.Handle("DELETE /api/items/", AuthMiddleware(tokens, http.HandlerFunc(h.DeleteItem)))

	cleanup := func() {
		_ = store.Close()
	}

	return h, tokens, mux, cleanup
}

func registerAndLogin(t *testing.T, router http.Handler, login string, password string) string {
	t.Helper()

	registerBody, _ := json.Marshal(dto.RegisterRequest{
		Login:    login,
		Password: password,
	})

	regReq := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(registerBody))
	regRec := httptest.NewRecorder()
	router.ServeHTTP(regRec, regReq)

	if regRec.Code != http.StatusCreated {
		t.Fatalf("unexpected register code: %d", regRec.Code)
	}

	var regResp dto.LoginResponse
	if err := json.NewDecoder(regRec.Body).Decode(&regResp); err != nil {
		t.Fatal(err)
	}

	return regResp.Token
}

func TestRegisterAndLogin(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")
	if token == "" {
		t.Fatal("expected token")
	}

	loginBody, _ := json.Marshal(dto.LoginRequest{
		Login:    "user1",
		Password: "pass1",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected login code: %d", rec.Code)
	}
}

func TestRegisterBadJSON(t *testing.T) {
	h, _, _, cleanup := newTestHTTP(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader([]byte("{bad")))
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestRegisterEmptyLogin(t *testing.T) {
	h, _, _, cleanup := newTestHTTP(t)
	defer cleanup()

	body, _ := json.Marshal(dto.RegisterRequest{
		Login:    "",
		Password: "pass1",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestRegisterConflict(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	_ = registerAndLogin(t, router, "user1", "pass1")

	body, _ := json.Marshal(dto.RegisterRequest{
		Login:    "user1",
		Password: "pass1",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestLoginBadPassword(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	_ = registerAndLogin(t, router, "user1", "pass1")

	body, _ := json.Marshal(dto.LoginRequest{
		Login:    "user1",
		Password: "wrong",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestItemCRUD(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")

	createBody, _ := json.Marshal(dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "telegram",
		Ciphertext: "cipher-1",
		Nonce:      "nonce-1",
		Salt:       "salt-1",
	})

	createReq := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusOK {
		t.Fatalf("unexpected create code: %d", createRec.Code)
	}

	var created dto.ItemResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if created.ID == "" {
		t.Fatal("expected item id")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected list code: %d", listRec.Code)
	}

	var items []dto.ItemResponse
	if err := json.NewDecoder(listRec.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/items/"+created.ID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get code: %d", getRec.Code)
	}

	var got dto.ItemResponse
	if err := json.NewDecoder(getRec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.Title != "note1" {
		t.Fatalf("unexpected title: %s", got.Title)
	}

	updateBody, _ := json.Marshal(dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note-updated",
		Meta:       "telegram-updated",
		Ciphertext: "cipher-2",
		Nonce:      "nonce-2",
		Salt:       "salt-2",
	})

	updateReq := httptest.NewRequest(http.MethodPut, "/api/items/"+created.ID, bytes.NewReader(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("unexpected update code: %d", updateRec.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/items/"+created.ID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("unexpected delete code: %d", deleteRec.Code)
	}
}

func TestItemUnauthorized(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestGetItemNotFound(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")

	req := httptest.NewRequest(http.MethodGet, "/api/items/missing", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestDeleteItemNotFound(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")

	req := httptest.NewRequest(http.MethodDelete, "/api/items/missing", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestUpsertItemBadBody(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}

func TestUpsertItemEmptyType(t *testing.T) {
	_, _, router, cleanup := newTestHTTP(t)
	defer cleanup()

	token := registerAndLogin(t, router, "user1", "pass1")

	body, _ := json.Marshal(dto.UpsertItemRequest{
		Type:       "",
		Title:      "note1",
		Meta:       "meta",
		Ciphertext: "cipher",
		Nonce:      "nonce",
		Salt:       "salt",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", rec.Code)
	}
}
