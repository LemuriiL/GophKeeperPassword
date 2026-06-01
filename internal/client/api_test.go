package client

import (
	"net/http/httptest"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

func TestAPIClientFlow(t *testing.T) {
	app, err := server.NewApp(server.Config{
		Address:   ":8080",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	ts := httptest.NewServer(app.RoutesForTests())
	defer ts.Close()

	api := NewAPIClient(ts.URL)

	token, salt, err := api.Register("user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("expected token")
	}

	if salt == "" {
		t.Fatal("expected salt")
	}

	token, err = api.Login("user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	item, err := api.SaveItem(token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "meta1",
		Ciphertext: "cipher",
		Nonce:      "nonce",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = api.GetItem(token, item.ID)
	if err != nil {
		t.Fatal(err)
	}

	items, err := api.ListItems(token)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	_, err = api.UpdateItem(token, item.ID, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note2",
		Meta:       "meta2",
		Ciphertext: "cipher2",
		Nonce:      "nonce2",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = api.DeleteItem(token, item.ID); err != nil {
		t.Fatal(err)
	}
}
