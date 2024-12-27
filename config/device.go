package config

import (
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/host"
	"strings"
)

var Device DeviceConfig

type DeviceConfig struct {
	ConfigurationURL string   `json:"configuration_url,omitempty"`
	Identifiers      []string `json:"identifiers,omitempty"`
	Name             string   `json:"name,omitempty"`
	Manufacturer     string   `json:"manufacturer,omitempty"`
	Model            string   `json:"model,omitempty"`
	SwVersion        string   `json:"sw_version,omitempty"`
	HwVersion        string   `json:"hw_version,omitempty"`
}

func loadDeviceConfig() {
	info, err := host.Info()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get host info")
	}

	if Config.HostID == "" {
		Config.HostID = info.HostID
	}
	if Config.DeviceName == "" {
		Config.DeviceName = info.Hostname
	}

	hwVersion, _, _ := strings.Cut(info.PlatformVersion, " Build ")

	Device = DeviceConfig{
		Identifiers:  []string{"system-bridge-" + Config.HostID},
		Name:         Config.DeviceName,
		Manufacturer: info.OS,
		Model:        info.Platform,
		SwVersion:    Version,
		HwVersion:    hwVersion,
	}
}
