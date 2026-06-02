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
}

func (s *stubServerApp) Run() error {
	return s.runErr
}

func (s *stubServerApp) Shutdown(ctx context.Context) error {
	return s.shutdownErr
}

func (s *stubServerApp) Close() error {
	return s.closeErr
}

func TestRunServerOK(t *testing.T) {
	oldLoad := loadServerConfig
	oldNew := newServerApp
	defer func() {
		loadServerConfig = oldLoad
		newServerApp = oldNew
	}()

	loadServerConfig = func(shortPath string, longPath string) (server.Config, error) {
		return server.Config{
			Address:   ":8080",
			DBPath:    ":memory:",
			JWTSecret: "secret",
		}, nil
	}

	newServerApp = func(cfg server.Config) (serverApp, error) {
		return &stubServerApp{}, nil
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-config", "server.json"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunServerLoadConfigError(t *testing.T) {
	oldLoad := loadServerConfig
	defer func() {
		loadServerConfig = oldLoad
	}()

	loadServerConfig = func(shortPath string, longPath string) (server.Config, error) {
		return server.Config{}, errors.New("bad config")
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-config", "server.json"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunServerNewAppError(t *testing.T) {
	oldLoad := loadServerConfig
	oldNew := newServerApp
	defer func() {
		loadServerConfig = oldLoad
		newServerApp = oldNew
	}()

	loadServerConfig = func(shortPath string, longPath string) (server.Config, error) {
		return server.Config{
			Address:   ":8080",
			DBPath:    ":memory:",
			JWTSecret: "secret",
		}, nil
	}

	newServerApp = func(cfg server.Config) (serverApp, error) {
		return nil, errors.New("new app failed")
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-config", "server.json"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunServerAppRunError(t *testing.T) {
	oldLoad := loadServerConfig
	oldNew := newServerApp
	defer func() {
		loadServerConfig = oldLoad
		newServerApp = oldNew
	}()

	loadServerConfig = func(shortPath string, longPath string) (server.Config, error) {
		return server.Config{
			Address:   ":8080",
			DBPath:    ":memory:",
			JWTSecret: "secret",
		}, nil
	}

	newServerApp = func(cfg server.Config) (serverApp, error) {
		return &stubServerApp{runErr: errors.New("run failed")}, nil
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-config", "server.json"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}
