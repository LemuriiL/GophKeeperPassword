package secure

import "testing"

func TestNewSalt(t *testing.T) {
	salt1, err := NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	salt2, err := NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	if salt1 == "" {
		t.Fatal("salt1 is empty")
	}

	if salt2 == "" {
		t.Fatal("salt2 is empty")
	}

	if salt1 == salt2 {
		t.Fatal("salts must be different")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	salt, err := NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	key := DeriveKey("master-password", []byte(salt))
	ciphertext, nonce, err := Encrypt(key, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	plain, err := Decrypt(key, ciphertext, nonce)
	if err != nil {
		t.Fatal(err)
	}

	if string(plain) != "hello world" {
		t.Fatalf("unexpected plaintext: %s", string(plain))
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1 := DeriveKey("master-password-1", []byte("salt-1"))
	key2 := DeriveKey("master-password-2", []byte("salt-2"))

	ciphertext, nonce, err := Encrypt(key1, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decrypt(key2, ciphertext, nonce)
	if err == nil {
		t.Fatal("expected decrypt error")
	}
}
