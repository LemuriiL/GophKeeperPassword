package buildinfo

import (
	"fmt"
	"io"
)

// Print печатает информацию о сборке
func Print(w io.Writer, version string, date string, commit string) {
	fmt.Fprintf(w, "Build version: %s\n", fallback(version))
	fmt.Fprintf(w, "Build date: %s\n", fallback(date))
	fmt.Fprintf(w, "Build commit: %s\n", fallback(commit))
}

// fallback возвращает значение по умолчанию
func fallback(v string) string {
	if v == "" {
		return "N/A"
	}

	return v
}
