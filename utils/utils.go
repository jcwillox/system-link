package utils

import (
	"os"
	"path/filepath"
)

func ExecutablePaths() (string, string, string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", "", "", err
	}
	return path, filepath.Dir(path), filepath.Base(path), nil
}
