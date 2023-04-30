package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/rs/zerolog/log"
)

var Config struct {
	MQTT struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		TLS            bool   `yaml:"tls"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		BaseTopic      string `yaml:"base_topic" env-default:"system-bridge"`
		DiscoveryTopic string `yaml:"discovery_topic" env-default:"homeassistant"`
	} `yaml:"mqtt"`
	Entities struct {
		Buttons       []buttons.Config        `yaml:"buttons"`
		Sensors       []sensors.Config        `yaml:"sensors"`
		BinarySensors []binary_sensors.Config `yaml:"binary_sensors"`
	} `yaml:"entities"`
	LogLevel     string `yaml:"log_level" env-default:"info"`
	LogLevelMqtt string `yaml:"log_level_mqtt" env-default:"error"`
}

func init() {
	setupLogging()

	// read config
	err := cleanenv.ReadConfig("config.yaml", &Config)
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading config")
	}

	setupLogLevels()

	log.Info().Msg("config loaded")

	setupDeviceConfig()
}
