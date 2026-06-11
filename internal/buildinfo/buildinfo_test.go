package buildinfo

import (
	"bytes"
	"strings"
	"testing"
)

// TestPrintFilledValues проверяет вывод заполненных build info значений
func TestPrintFilledValues(t *testing.T) {
	var out bytes.Buffer

	Print(&out, "1.0.0", "2026-06-01", "abc123")

	got := out.String()

	if !strings.Contains(got, "Build version: 1.0.0") {
		t.Fatalf("unexpected output: %s", got)
	}

	if !strings.Contains(got, "Build date: 2026-06-01") {
		t.Fatalf("unexpected output: %s", got)
	}

	if !strings.Contains(got, "Build commit: abc123") {
		t.Fatalf("unexpected output: %s", got)
	}
}

// TestPrintEmptyValues проверяет вывод значений по умолчанию
func TestPrintEmptyValues(t *testing.T) {
	var out bytes.Buffer

	Print(&out, "", "", "")

	got := out.String()

	if strings.Count(got, "N/A") != 3 {
		t.Fatalf("unexpected output: %s", got)
	}
}
