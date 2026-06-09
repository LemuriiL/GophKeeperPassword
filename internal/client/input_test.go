package client

import "testing"

func TestReadSecretFromStdin(t *testing.T) {
	withStdin(t, "hello\n", func() {
		value, err := ReadSecret("Password: ")
		if err != nil {
			t.Fatal(err)
		}

		if value != "hello" {
			t.Fatalf("unexpected value: %s", value)
		}
	})
}

func TestReadSecretTrimsSpaces(t *testing.T) {
	withStdin(t, "  hello  \n", func() {
		value, err := ReadSecret("Password: ")
		if err != nil {
			t.Fatal(err)
		}

		if value != "hello" {
			t.Fatalf("unexpected value: %s", value)
		}
	})
}
