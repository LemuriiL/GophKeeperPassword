package client

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadSecret читает секрет без эха
func ReadSecret(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)

	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}

func resolveSecret(flagValue string, prompt string) (string, error) {
	if strings.TrimSpace(flagValue) != "" {
		return flagValue, nil
	}

	return ReadSecret(prompt)
}
