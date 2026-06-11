package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

// newTestHTTPServer создает тестовый HTTP сервер
func newTestHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	app, err := NewApp(Config{
		Address:     ":0",
		DBPath:      ":memory:",
		JWTSecret:   "test-secret",
		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",
	})
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(app.RoutesForTests())

	cleanup := func() {
		ts.Close()
		_ = app.Close()
	}

	return ts, cleanup
}

// doJSON выполняет JSON запрос
func doJSON(t *testing.T, method string, url string, token string, body any) *http.Response {
	t.Helper()

	var reader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

// readBody читает тело ответа
func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	return string(data)
}

// registerUser регистрирует пользователя и возвращает токен
func registerUser(t *testing.T, ts *httptest.Server, login string, password string) string {
	t.Helper()

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/register", "", dto.RegisterRequest{
		Login:    login,
		Password: password,
	})

	if resp.StatusCode != http.StatusCreated {
		body := readBody(t, resp)
		t.Fatalf("unexpected register code: %d body: %s", resp.StatusCode, body)
	}

	var out dto.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(out.Token) == "" {
		t.Fatal("expected token")
	}

	return out.Token
}

// TestRegister проверяет регистрацию пользователя
func TestRegister(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/register", "", dto.RegisterRequest{
		Login:    "user1",
		Password: "pass1",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}

	var out dto.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(out.Token) == "" {
		t.Fatal("expected token")
	}
}

// TestRegisterDuplicate проверяет конфликт при повторной регистрации
func TestRegisterDuplicate(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	_ = registerUser(t, ts, "user1", "pass1")

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/register", "", dto.RegisterRequest{
		Login:    "user1",
		Password: "pass1",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}

// TestRegisterBadRequest проверяет ошибку регистрации с пустыми данными
func TestRegisterBadRequest(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/register", "", dto.RegisterRequest{
		Login:    "",
		Password: "",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}

// TestLogin проверяет логин пользователя
func TestLogin(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	_ = registerUser(t, ts, "user1", "pass1")

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/login", "", dto.LoginRequest{
		Login:    "user1",
		Password: "pass1",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}

	var out dto.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(out.Token) == "" {
		t.Fatal("expected token")
	}
}

// TestLoginInvalidCredentials проверяет общий ответ при неверном логине
func TestLoginInvalidCredentials(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/login", "", dto.LoginRequest{
		Login:    "missing",
		Password: "pass1",
	})

	if resp.StatusCode != http.StatusUnauthorized {
		body := readBody(t, resp)
		t.Fatalf("unexpected code: %d body: %s", resp.StatusCode, body)
	}

	body := readBody(t, resp)

	if strings.Contains(strings.ToLower(body), "sql") {
		t.Fatalf("response leaks internal error: %s", body)
	}

	if !strings.Contains(body, "invalid login or password") {
		t.Fatalf("unexpected body: %s", body)
	}
}

// TestItemCRUD проверяет создание, чтение, обновление и удаление секрета
func TestItemCRUD(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	token := registerUser(t, ts, "user1", "pass1")

	createResp := doJSON(t, http.MethodPost, ts.URL+"/api/items", token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "test",
		Ciphertext: "ciphertext1",
		Nonce:      "nonce1",
		Salt:       "salt1",
	})
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected create code: %d", createResp.StatusCode)
	}

	var created dto.ItemResponse
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(created.ID) == "" {
		t.Fatal("expected item id")
	}

	if created.Title != "note1" {
		t.Fatalf("unexpected title: %s", created.Title)
	}

	listResp := doJSON(t, http.MethodGet, ts.URL+"/api/items", token, nil)
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected list code: %d", listResp.StatusCode)
	}

	var list []dto.ItemResponse
	if err := json.NewDecoder(listResp.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Fatalf("unexpected list len: %d", len(list))
	}

	getResp := doJSON(t, http.MethodGet, ts.URL+"/api/items/"+created.ID, token, nil)
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected get code: %d", getResp.StatusCode)
	}

	var got dto.ItemResponse
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.ID != created.ID {
		t.Fatalf("unexpected id: %s", got.ID)
	}

	updateResp := doJSON(t, http.MethodPut, ts.URL+"/api/items/"+created.ID, token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note2",
		Meta:       "updated",
		Ciphertext: "ciphertext2",
		Nonce:      "nonce2",
		Salt:       "salt2",
	})
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected update code: %d", updateResp.StatusCode)
	}

	var updated dto.ItemResponse
	if err := json.NewDecoder(updateResp.Body).Decode(&updated); err != nil {
		t.Fatal(err)
	}

	if updated.Title != "note2" {
		t.Fatalf("unexpected updated title: %s", updated.Title)
	}

	if updated.Ciphertext != "ciphertext2" {
		t.Fatalf("unexpected updated ciphertext: %s", updated.Ciphertext)
	}

	deleteResp := doJSON(t, http.MethodDelete, ts.URL+"/api/items/"+created.ID, token, nil)
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected delete code: %d", deleteResp.StatusCode)
	}

	getDeletedResp := doJSON(t, http.MethodGet, ts.URL+"/api/items/"+created.ID, token, nil)
	defer getDeletedResp.Body.Close()

	if getDeletedResp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected get deleted code: %d", getDeletedResp.StatusCode)
	}
}

// TestCreateItemBadRequest проверяет создание секрета с некорректными данными
func TestCreateItemBadRequest(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	token := registerUser(t, ts, "user1", "pass1")

	resp := doJSON(t, http.MethodPost, ts.URL+"/api/items", token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "test",
		Ciphertext: "",
		Nonce:      "",
		Salt:       "",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}

// TestUpdateMissingItem проверяет обновление отсутствующего секрета
func TestUpdateMissingItem(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	token := registerUser(t, ts, "user1", "pass1")

	resp := doJSON(t, http.MethodPut, ts.URL+"/api/items/missing-id", token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "test",
		Ciphertext: "ciphertext1",
		Nonce:      "nonce1",
		Salt:       "salt1",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}

// TestDeleteMissingItem проверяет удаление отсутствующего секрета
func TestDeleteMissingItem(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	token := registerUser(t, ts, "user1", "pass1")

	resp := doJSON(t, http.MethodDelete, ts.URL+"/api/items/missing-id", token, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}

// TestUnauthorizedItems проверяет запрет доступа без токена
func TestUnauthorizedItems(t *testing.T) {
	ts, cleanup := newTestHTTPServer(t)
	defer cleanup()

	resp := doJSON(t, http.MethodGet, ts.URL+"/api/items", "", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected code: %d", resp.StatusCode)
	}
}
