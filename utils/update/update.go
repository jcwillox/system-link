package update

import (
	"fmt"
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/utils"
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
	exePath, exeDir, err := utils.ExecutablePaths()
	if err != nil {
		return err
	}
	log.Debug().Str("dir", exeDir).Msg("executable directory")

	// download file
	archivePath := filepath.Join(exeDir, "update.tar.gz")
	if runtime.GOOS == "windows" {
		archivePath = filepath.Join(exeDir, "update.zip")
	}

	err = downloadFile(downloadUrl, archivePath)
	if err != nil {
		return err
	}
	log.Info().Str("latest_version", latestVersion).Msg("update downloaded")

	// extract file
	outputPath := exePath + ".new"
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
	err = os.Rename(exePath, exePath+".old")
	if err != nil {
		return fmt.Errorf("failed to rename executable to '.old': %w", err)
	}

	// rename new executable to current
	err = os.Rename(outputPath, exePath)
	if err != nil {
		return fmt.Errorf("failed to rename executable to '.new': %w", err)
	}

	// delete old executable
	// ignore error as it will fail on windows as executable is still running
	_ = os.Remove(exePath + ".old")

	log.Info().Msg("restarting system-bridge for update")

	if runtime.GOOS == "windows" {
		_, err := os.StartProcess(exePath, os.Args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if err != nil {
			log.Err(err).Msg("failed to automatically restart system-bridge")
			return err
		}
		os.Exit(0)
	} else {
		// exec is preferred over StartProcess as it will replace the current process,
		// but it is not available on windows
		err := syscall.Exec(exePath, os.Args, os.Environ())
		if err != nil {
			log.Err(err).Msg("failed to automatically restart system-bridge")
			return err
		}
	}

	return nil
}

func Cleanup() {
	exePath, _, err := utils.ExecutablePaths()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get executable path")
	}
	_ = os.Remove(exePath + ".old")
}
