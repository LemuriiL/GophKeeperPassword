package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

func newTestClientApp(t *testing.T) (*App, func()) {
	t.Helper()

	srvApp, err := server.NewApp(server.Config{
		Address:   ":0",
		DBPath:    ":memory:",
		JWTSecret: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(srvApp.RoutesForTests())

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "client.json")
	sessionPath := filepath.Join(dir, ".gophkeeper_session.json")

	cfgData, err := json.Marshal(map[string]any{
		"server_url":   ts.URL,
		"session_file": sessionPath,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile(cfgPath, cfgData, 0o600); err != nil {
		t.Fatal(err)
	}

	app, err := NewApp(cfgPath, "")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		ts.Close()
		_ = srvApp.Close()
	}

	return app, cleanup
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}

func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()

	old := os.Stdin

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.WriteString(input)
	if err != nil {
		t.Fatal(err)
	}

	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	os.Stdin = r
	defer func() {
		os.Stdin = old
		_ = r.Close()
	}()

	fn()
}

func TestAppRunRegisterLoginAddListGetUpdateDelete(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	withStdin(t, "pass1\n", func() {
		if err := app.Run([]string{"register", "--login", "user1"}); err != nil {
			t.Fatal(err)
		}
	})

	withStdin(t, "pass1\n", func() {
		if err := app.Run([]string{"login", "--login", "user1"}); err != nil {
			t.Fatal(err)
		}
	})

	session, err := LoadSession(app.cfg.SessionFile)
	if err != nil {
		t.Fatal(err)
	}

	if session.Login != "user1" {
		t.Fatalf("unexpected login: %s", session.Login)
	}

	if strings.TrimSpace(session.Token) == "" {
		t.Fatal("expected token after login")
	}

	var itemID string
	out := captureStdout(t, func() {
		withStdin(t, "local-master\n", func() {
			err := app.Run([]string{
				"add",
				"--type", "text",
				"--title", "note1",
				"--meta", "test",
				"--value", `{"text":"hello"}`,
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})
	itemID = strings.TrimSpace(out)
	if itemID == "" {
		t.Fatal("expected item id")
	}

	out = captureStdout(t, func() {
		withStdin(t, "local-master\n", func() {
			err := app.Run([]string{
				"list",
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	if !strings.Contains(out, "Title: note1") {
		t.Fatalf("unexpected list output: %s", out)
	}

	if !strings.Contains(out, `{"text":"hello"}`) {
		t.Fatalf("unexpected list output: %s", out)
	}

	out = captureStdout(t, func() {
		withStdin(t, "local-master\n", func() {
			err := app.Run([]string{
				"get",
				"--id", itemID,
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	if !strings.Contains(out, "ID: "+itemID) {
		t.Fatalf("unexpected get output: %s", out)
	}

	withStdin(t, "local-master\n", func() {
		if err := app.Run([]string{
			"update",
			"--id", itemID,
			"--type", "text",
			"--title", "note2",
			"--meta", "test2",
			"--value", `{"text":"updated"}`,
		}); err != nil {
			t.Fatal(err)
		}
	})

	out = captureStdout(t, func() {
		withStdin(t, "local-master\n", func() {
			err := app.Run([]string{
				"get",
				"--id", itemID,
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	if !strings.Contains(out, "Title: note2") {
		t.Fatalf("unexpected updated get output: %s", out)
	}

	if !strings.Contains(out, `{"text":"updated"}`) {
		t.Fatalf("unexpected updated get output: %s", out)
	}

	if err := app.Run([]string{
		"delete",
		"--id", itemID,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestAppRunUnknownCommand(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	err := app.Run([]string{"bad-command"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAppRunNoCommand(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	err := app.Run(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAppRunListWithoutSession(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	withStdin(t, "local-master\n", func() {
		err := app.Run([]string{
			"list",
		})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestAppRunGetWithoutSession(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	withStdin(t, "local-master\n", func() {
		err := app.Run([]string{
			"get",
			"--id", "missing",
		})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestAppRunDeleteWithoutSession(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	err := app.Run([]string{
		"delete",
		"--id", "missing",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAppRegisterCreatesSession(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	withStdin(t, "pass3\n", func() {
		if err := app.Run([]string{"register", "--login", "user3"}); err != nil {
			t.Fatal(err)
		}
	})

	session, err := LoadSession(app.cfg.SessionFile)
	if err != nil {
		t.Fatal(err)
	}

	if session.Login != "user3" {
		t.Fatalf("unexpected login: %s", session.Login)
	}

	if strings.TrimSpace(session.Token) == "" {
		t.Fatal("expected token")
	}
}

func TestAppLoginWritesSession(t *testing.T) {
	app, cleanup := newTestClientApp(t)
	defer cleanup()

	withStdin(t, "pass2\n", func() {
		if err := app.Run([]string{"register", "--login", "user2"}); err != nil {
			t.Fatal(err)
		}
	})

	withStdin(t, "pass2\n", func() {
		if err := app.Run([]string{"login", "--login", "user2"}); err != nil {
			t.Fatal(err)
		}
	})

	session, err := LoadSession(app.cfg.SessionFile)
	if err != nil {
		t.Fatal(err)
	}

	if session.Login != "user2" {
		t.Fatalf("unexpected login: %s", session.Login)
	}

	if strings.TrimSpace(session.Token) == "" {
		t.Fatal("expected token after login")
	}
}
