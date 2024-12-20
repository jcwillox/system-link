package config

import (
	_ "embed"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"path"
)

var (
	RepoUrl         = "https://github.com/jcwillox/system-bridge"
	Version         string
	Config          CoreConfig
	ShutdownChannel = make(chan bool)
)

type CoreConfig struct {
	MQTT struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		TLS            bool   `yaml:"tls"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		BaseTopic      string `yaml:"base_topic" env-default:"system-bridge"`
		DiscoveryTopic string `yaml:"discovery_topic" env-default:"homeassistant"`
	} `yaml:"mqtt"`
	LogLevel     string `yaml:"log_level" env-default:"info"`
	LogLevelMqtt string `yaml:"log_level_mqtt" env-default:"error"`
}

func (c *CoreConfig) AvailabilityTopic() string {
	return path.Join(Config.MQTT.BaseTopic, HostID, "availability")
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
