package update

import (
	"github.com/jcwillox/system-link/config"
	"net/http"
	"path"
	"runtime"
	"strings"
)

func GetLatestVersion() (string, error) {
	resp, err := http.Head(config.RepoUrl + "/releases/latest")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(path.Base(resp.Request.URL.Path), "v"), nil
}

func GetDownloadURL(version string) string {
	downloadUrl := config.RepoUrl + "/releases/download/v" + version +
		"/system_link_" + version + "_" + runtime.GOOS + "_" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		return downloadUrl + ".zip"
	}
	return downloadUrl + ".tar.gz"
}
