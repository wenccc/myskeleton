package file

import (
	"os"
)

func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Put(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}
