package utils

import (
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	ExeName      string
	ExeBaseName  string
	ExePath      string
	ExePathHash  string
	ExeDirectory string
	UserHomeDir  = sync.OnceValue(func() string {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get user home directory")
		}
		return homeDir
	})
	// ShareDirectory is the directory where the app stores its data in non-portable mode
	ShareDirectory = sync.OnceValue(func() string {
		return filepath.Join(UserHomeDir(), ".local/share/system-link")
	})
	PortableMode = sync.OnceValue(func() bool {
		// if config exists in the exe directory, we are in portable mode
		return FileExists(filepath.Join(ExeDirectory, "config.yaml"))
	})
	ConfigPath = sync.OnceValue(func() string {
		// load from exe directory in portable mode
		if PortableMode() {
			return filepath.Join(ExeDirectory, "config.yaml")
		}
		// try to load from the env path
		if env := os.Getenv("SYSTEM_LINK_CONFIG"); env != "" {
			return env
		}
		// try load from ~/.config
		return filepath.Join(UserHomeDir(), ".config/system-link/config.yaml")
	})
	ConfigDirectory = sync.OnceValue(func() string {
		return filepath.Dir(ConfigPath())
	})
	LogsPath = sync.OnceValue(func() string {
		// use exe directory if in portable mode
		if PortableMode() {
			return filepath.Join(ExeDirectory, ExeBaseName+".log")
		}
		logName := ExeBaseName + "." + ExePathHash + ".log"
		// use env path if set
		if env := os.Getenv("SYSTEM_LINK_LOGS_DIR"); env != "" {
			return filepath.Join(env, logName)
		}
		// otherwise use the share directory
		return filepath.Join(ShareDirectory(), logName)
	})
	LockPath = sync.OnceValue(func() string {
		if PortableMode() {
			return ExePath + ".lock"
		}
		return filepath.Join(ShareDirectory(), ExeBaseName+"."+ExePathHash+".lock")
	})
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Fnv1aHash(s string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 36)
}

func init() {
	path, err := os.Executable()
	if err != nil {
		log.Panic().Err(err).Msg("failed to get executable path")
	}
	ExePath = path
	ExeDirectory = filepath.Dir(path)
	ExeName = filepath.Base(path)
	ExeBaseName = strings.TrimSuffix(ExeName, filepath.Ext(ExeName))
	ExePathHash = Fnv1aHash(ExePath)
}
