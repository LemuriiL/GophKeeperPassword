package client

import (
	"net/http/httptest"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/model"
	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
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

	token, _, err := api.Register("user1", "pass1")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("expected token")
	}

	itemSalt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := `{"text":"hello"}`
	ciphertext, nonce, err := EncryptPayload("local-master", itemSalt, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	item, err := api.SaveItem(token, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note1",
		Meta:       "meta1",
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       itemSalt,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := api.GetItem(token, item.ID)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := DecryptPayload("local-master", got.Salt, got.Ciphertext, got.Nonce)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Fatalf("unexpected decrypted payload: %s", decrypted)
	}

	items, err := api.ListItems(token)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	decrypted, err = DecryptPayload("local-master", items[0].Salt, items[0].Ciphertext, items[0].Nonce)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Fatalf("unexpected decrypted payload from list: %s", decrypted)
	}

	newItemSalt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	newPlaintext := `{"text":"updated"}`
	newCiphertext, newNonce, err := EncryptPayload("local-master", newItemSalt, newPlaintext)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := api.UpdateItem(token, item.ID, dto.UpsertItemRequest{
		Type:       model.TypeText,
		Title:      "note2",
		Meta:       "meta2",
		Ciphertext: newCiphertext,
		Nonce:      newNonce,
		Salt:       newItemSalt,
	})
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err = DecryptPayload("local-master", updated.Salt, updated.Ciphertext, updated.Nonce)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != newPlaintext {
		t.Fatalf("unexpected updated decrypted payload: %s", decrypted)
	}

	if err = api.DeleteItem(token, item.ID); err != nil {
		t.Fatal(err)
	}
}
