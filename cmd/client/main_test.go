package main

import (
	"bytes"
	"errors"
	"testing"
)

type stubClientApp struct {
	runErr  error
	gotArgs []string
}

func (s *stubClientApp) Run(args []string) error {
	s.gotArgs = append([]string(nil), args...)
	return s.runErr
}

func TestRunClientOK(t *testing.T) {
	oldFactory := newClientApp
	defer func() {
		newClientApp = oldFactory
	}()

	app := &stubClientApp{}

	newClientApp = func(configShort string, configLong string) (clientApp, error) {
		if configShort != "" {
			t.Fatalf("unexpected configShort: %s", configShort)
		}
		if configLong != "client.json" {
			t.Fatalf("unexpected configLong: %s", configLong)
		}
		return app, nil
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-config", "client.json", "list"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("unexpected code: %d", code)
	}

	if len(app.gotArgs) == 0 || app.gotArgs[0] != "list" {
		t.Fatalf("unexpected args: %#v", app.gotArgs)
	}
}

func TestRunClientFactoryError(t *testing.T) {
	oldFactory := newClientApp
	defer func() {
		newClientApp = oldFactory
	}()

	newClientApp = func(configShort string, configLong string) (clientApp, error) {
		return nil, errors.New("boom")
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"list"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunClientAppError(t *testing.T) {
	oldFactory := newClientApp
	defer func() {
		newClientApp = oldFactory
	}()

	app := &stubClientApp{runErr: errors.New("run failed")}

	newClientApp = func(configShort string, configLong string) (clientApp, error) {
		return app, nil
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"list"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunClientFlagError(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"-bad-flag"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("unexpected code: %d", code)
	}
}
