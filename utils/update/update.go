package update

import (
	"fmt"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"runtime"
)

func Update(updateProgress func(progress float64)) error {
	updateProgress(20)

	// get latest version
	latestVersion, err := GetLatestVersion()
	log.Debug().Str("latest_version", latestVersion).Msg("latest version")
	if err != nil || latestVersion == config.Version {
		return err
	}
	updateProgress(40)

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
	updateProgress(60)

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
	updateProgress(80)

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
	updateProgress(100)

	log.Info().Msg("restarting system-link for update")
	return utils.RestartSelf()
}

func Cleanup() {
	_ = os.Remove(utils.ExePath + ".old")
}
