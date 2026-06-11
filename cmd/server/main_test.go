package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

type stubServerApp struct {
	runErr      error
	shutdownErr error
	closeErr    error
	closed      bool
}

func (s *stubServerApp) Run() error {
	return s.runErr
}

func (s *stubServerApp) Shutdown(ctx context.Context) error {
	return s.shutdownErr
}

func (s *stubServerApp) Close() error {
	s.closed = true
	return s.closeErr
}

func TestRunServerVersionFlag(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("unexpected code: %d", code)
	}

	if !bytes.Contains(out.Bytes(), []byte("Build version:")) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRunServerVersionCommand(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("unexpected code: %d", code)
	}

	if !bytes.Contains(out.Bytes(), []byte("Build version:")) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRunServerConfigError(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}

	if !bytes.Contains(errOut.Bytes(), []byte("JWT_SECRET is required")) {
		t.Fatalf("unexpected stderr: %s", errOut.String())
	}
}

func TestRunServerAppFactoryError(t *testing.T) {
	oldFactory := newServerApp
	defer func() {
		newServerApp = oldFactory
	}()

	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("TLS_CERT_FILE", "cert.pem")
	t.Setenv("TLS_KEY_FILE", "key.pem")

	newServerApp = func(cfg server.Config) (serverApp, error) {
		return nil, errors.New("factory failed")
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}

	if !bytes.Contains(errOut.Bytes(), []byte("factory failed")) {
		t.Fatalf("unexpected stderr: %s", errOut.String())
	}
}

func TestRunServerRunError(t *testing.T) {
	oldFactory := newServerApp
	defer func() {
		newServerApp = oldFactory
	}()

	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("TLS_CERT_FILE", "cert.pem")
	t.Setenv("TLS_KEY_FILE", "key.pem")

	app := &stubServerApp{
		runErr: errors.New("run failed"),
	}

	newServerApp = func(cfg server.Config) (serverApp, error) {
		return app, nil
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}

	if !app.closed {
		t.Fatal("expected app close")
	}

	if !bytes.Contains(errOut.Bytes(), []byte("run failed")) {
		t.Fatalf("unexpected stderr: %s", errOut.String())
	}
}

func TestRunServerFlagError(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-bad-flag"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}
