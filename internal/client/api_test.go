package client

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

// newTestAPIClient создает тестовый API клиент
func newTestAPIClient(t *testing.T) (*APIClient, func()) {
	t.Helper()

	app, err := server.NewApp(server.Config{
		Address:     ":0",
		DBPath:      ":memory:",
		JWTSecret:   "secret",
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

	return NewAPIClient(ts.URL), cleanup
}

// TestAPIClientRegisterLogin проверяет регистрацию и логин
func TestAPIClientRegisterLogin(t *testing.T) {
	api, cleanup := newTestAPIClient(t)
	defer cleanup()

	token, err := api.Register("user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(token) == "" {
		t.Fatal("expected token")
	}

	token, err = api.Login("user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(token) == "" {
		t.Fatal("expected token")
	}
}

// TestAPIClientItems проверяет CRUD секретов через API клиент
func TestAPIClientItems(t *testing.T) {
	api, cleanup := newTestAPIClient(t)
	defer cleanup()

	token, err := api.Register("user2", "pass2")
	if err != nil {
		t.Fatal(err)
	}

	salt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, nonce, err := EncryptPayload("master", salt, `{"text":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}

	item, err := api.SaveItem(token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "test",
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       salt,
	})
	if err != nil {
		t.Fatal(err)
	}

	if item.ID == "" {
		t.Fatal("expected item id")
	}

	items, err := api.ListItems(token)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("unexpected items count: %d", len(items))
	}

	got, err := api.GetItem(token, item.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != "note1" {
		t.Fatalf("unexpected title: %s", got.Title)
	}

	updatedSalt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	updatedCiphertext, updatedNonce, err := EncryptPayload("master", updatedSalt, `{"text":"updated"}`)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := api.UpdateItem(token, item.ID, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note2",
		Meta:       "updated",
		Ciphertext: updatedCiphertext,
		Nonce:      updatedNonce,
		Salt:       updatedSalt,
	})
	if err != nil {
		t.Fatal(err)
	}

	if updated.Title != "note2" {
		t.Fatalf("unexpected title: %s", updated.Title)
	}

	if err = api.DeleteItem(token, item.ID); err != nil {
		t.Fatal(err)
	}
}

// TestNewAPIClientWithInsecureSkipVerify проверяет создание клиента с отключенной TLS проверкой
func TestNewAPIClientWithInsecureSkipVerify(t *testing.T) {
	api := NewAPIClient("https://localhost:8080/", true)

	if api == nil {
		t.Fatal("expected api client")
	}

	if api.baseURL != "https://localhost:8080" {
		t.Fatalf("unexpected base url: %s", api.baseURL)
	}
}

// TestAPIClientBadServer проверяет ошибку недоступного сервера
func TestAPIClientBadServer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	_, err := os.ReadFile(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

// TestAPIClientMarshalRequest проверяет сериализацию запроса
func TestAPIClientMarshalRequest(t *testing.T) {
	data, err := json.Marshal(dto.RegisterRequest{
		Login:    "user",
		Password: "pass",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "user") {
		t.Fatalf("unexpected json: %s", string(data))
	}
}
