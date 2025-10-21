package utils

import "os"

func IsSystemd() bool {
	if _, err := os.Stat("/run/systemd/system"); !os.IsNotExist(err) {
		return true
	}
	return false
}
