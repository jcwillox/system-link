package config

import (
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/host"
	"strings"
)

var Device DeviceConfig
var HostID string

type DeviceConfig struct {
	ConfigurationURL string   `json:"configuration_url,omitempty"`
	Identifiers      []string `json:"identifiers,omitempty"`
	Name             string   `json:"name,omitempty"`
	Manufacturer     string   `json:"manufacturer,omitempty"`
	Model            string   `json:"model,omitempty"`
	SwVersion        string   `json:"sw_version,omitempty"`
	HwVersion        string   `json:"hw_version,omitempty"`
}

func setupDeviceConfig() {
	info, err := host.Info()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get host info")
	}

	hwVersion, _, _ := strings.Cut(info.PlatformVersion, " Build ")

	HostID = info.HostID
	Device = DeviceConfig{
		Identifiers:  []string{"system-bridge", info.HostID},
		Name:         info.Hostname,
		Manufacturer: info.OS,
		Model:        info.Platform,
		SwVersion:    Version,
		HwVersion:    hwVersion,
	}
}
