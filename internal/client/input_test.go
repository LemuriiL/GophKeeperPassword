package client

import "testing"

func TestResolveSecretUsesFlagValue(t *testing.T) {
	value, err := resolveSecret("hello", "Password: ")
	if err != nil {
		t.Fatal(err)
	}

	if value != "hello" {
		t.Fatalf("unexpected value: %s", value)
	}
}
