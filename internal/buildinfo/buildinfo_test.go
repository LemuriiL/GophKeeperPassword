package buildinfo

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}

func TestPrint(t *testing.T) {
	out := captureOutput(t, func() {
		Print("1.0.0", "2026-06-02", "abc123")
	})

	if !strings.Contains(out, "Build version: 1.0.0") {
		t.Fatalf("unexpected output: %s", out)
	}

	if !strings.Contains(out, "Build date: 2026-06-02") {
		t.Fatalf("unexpected output: %s", out)
	}

	if !strings.Contains(out, "Build commit: abc123") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestPrintFallback(t *testing.T) {
	out := captureOutput(t, func() {
		Print("", "", "")
	})

	if !strings.Contains(out, "Build version: N/A") {
		t.Fatalf("unexpected output: %s", out)
	}

	if !strings.Contains(out, "Build date: N/A") {
		t.Fatalf("unexpected output: %s", out)
	}

	if !strings.Contains(out, "Build commit: N/A") {
		t.Fatalf("unexpected output: %s", out)
	}
}
