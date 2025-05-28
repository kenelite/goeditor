package backend

import (
	"os"
)

func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func SaveFile(path, content string) {
	_ = os.WriteFile(path, []byte(content), 0644)
}
