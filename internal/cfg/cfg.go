package cfg

import (
	"encoding/json"
	"os"
)

// Load читает JSON-конфиг из файла.
func Load[T any](path string) (T, error) {
	var out T

	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(data, &out)
	return out, err
}
