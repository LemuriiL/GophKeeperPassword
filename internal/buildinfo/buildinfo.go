package buildinfo

import "fmt"

// Print печатает информацию о сборке.
func Print(version string, date string, commit string) {
	fmt.Printf("Build version: %s\n", fallback(version))
	fmt.Printf("Build date: %s\n", fallback(date))
	fmt.Printf("Build commit: %s\n", fallback(commit))
}

func fallback(v string) string {
	if v == "" {
		return "N/A"
	}

	return v
}
