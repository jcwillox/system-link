package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"syscall"
)

func RestartSelf() error {
	if runtime.GOOS == "windows" {
		_, err := os.StartProcess(ExePath, os.Args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if err != nil {
			log.Err(err).Msg("failed to automatically restart system-link")
			return err
		}
		os.Exit(0)
	} else {
		// exec is preferred over StartProcess as it will replace the current process,
		// but it is not available on windows
		err := syscall.Exec(ExePath, os.Args, os.Environ())
		if err != nil {
			log.Err(err).Msg("failed to automatically restart system-link")
			return err
		}
	}
	return nil
}
