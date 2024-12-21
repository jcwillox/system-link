package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
	"os"
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
		BaseTopic      string `yaml:"base_topic" default:"system-bridge"`
		DiscoveryTopic string `yaml:"discovery_topic" default:"homeassistant"`
	} `yaml:"mqtt"`
	LogLevel     string `yaml:"log_level" default:"info"`
	LogLevelMqtt string `yaml:"log_level_mqtt" default:"error"`
}

func (c *CoreConfig) AvailabilityTopic() string {
	return path.Join(Config.MQTT.BaseTopic, HostID, "availability")
}

func LoadConfig() {
	// set defaults
	err := defaults.Set(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set config defaults")
	}

	// load config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading config")
	}

	// parse config
	if err = yaml.Unmarshal(data, &Config); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing config")
	}

	log.Info().Msg("config loaded")
}
