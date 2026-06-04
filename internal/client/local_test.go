package client

import (
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
)

func TestEncryptPayloadDecryptPayload(t *testing.T) {
	salt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, nonce, err := EncryptPayload("master-password", salt, `{"text":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}

	plain, err := DecryptPayload("master-password", salt, ciphertext, nonce)
	if err != nil {
		t.Fatal(err)
	}

	if plain != `{"text":"hello"}` {
		t.Fatalf("unexpected plaintext: %s", plain)
	}
}

func TestDecryptPayloadWrongPassword(t *testing.T) {
	salt, err := secure.NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, nonce, err := EncryptPayload("master-password", salt, `{"text":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = DecryptPayload("wrong-password", salt, ciphertext, nonce)
	if err == nil {
		t.Fatal("expected decrypt error")
	}
}
