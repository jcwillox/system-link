package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

var (
	ExeName      string
	ExePath      string
	ExeDirectory string
)

func init() {
	path, err := os.Executable()
	if err != nil {
		log.Panic().Err(err).Msg("failed to get executable path")
	}
	ExePath = path
	ExeDirectory = filepath.Dir(path)
	ExeName = filepath.Base(path)
}
