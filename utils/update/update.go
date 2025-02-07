package update

import (
	"fmt"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func Update() error {
	// get latest version
	latestVersion, err := GetLatestVersion()
	log.Debug().Str("latest_version", latestVersion).Msg("latest version")
	if err != nil || latestVersion == config.Version {
		return err
	}

	// get download url
	downloadUrl := GetDownloadURL(latestVersion)
	log.Debug().Str("download_url", downloadUrl).Msg("downloading update")

	// get executable directory (working directory)
	log.Debug().Str("dir", utils.ExeDirectory).Msg("executable directory")

	// download file
	archivePath := filepath.Join(utils.ExeDirectory, "update.tar.gz")
	if runtime.GOOS == "windows" {
		archivePath = filepath.Join(utils.ExeDirectory, "update.zip")
	}

	err = downloadFile(downloadUrl, archivePath)
	if err != nil {
		return err
	}
	log.Info().Str("latest_version", latestVersion).Msg("update downloaded")

	// extract file
	outputPath := utils.ExePath + ".new"
	if runtime.GOOS == "windows" {
		err = extractZip(archivePath, outputPath)
		if err != nil {
			return err
		}
	} else {
		err = extractTarGz(archivePath, outputPath)
		if err != nil {
			return err
		}
		err := os.Chmod(outputPath, 0755)
		if err != nil {
			return err
		}
	}

	// delete archive
	err = os.Remove(archivePath)
	if err != nil {
		return err
	}

	// rename current executable to old
	err = os.Rename(utils.ExePath, utils.ExePath+".old")
	if err != nil {
		return fmt.Errorf("failed to rename executable to '.old': %w", err)
	}

	// rename new executable to current
	err = os.Rename(outputPath, utils.ExePath)
	if err != nil {
		return fmt.Errorf("failed to rename executable to '.new': %w", err)
	}

	// delete old executable
	// ignore error as it will fail on windows as executable is still running
	_ = os.Remove(utils.ExePath + ".old")

	log.Info().Msg("restarting system-link for update")

	if runtime.GOOS == "windows" {
		_, err := os.StartProcess(utils.ExePath, os.Args, &os.ProcAttr{
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
		err := syscall.Exec(utils.ExePath, os.Args, os.Environ())
		if err != nil {
			log.Err(err).Msg("failed to automatically restart system-link")
			return err
		}
	}

	return nil
}

func Cleanup() {
	_ = os.Remove(utils.ExePath + ".old")
}
