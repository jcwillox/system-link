package utils

import (
	"os"
	"path/filepath"
)

func ExecutablePaths() (string, string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", "", err
	}
	return path, filepath.Dir(path), nil
}
