package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
)

var (
	RepoUrl         = "https://github.com/jcwillox/system-bridge"
	Version         string
	Config          CoreConfig
	ShutdownChannel = make(chan bool)
)

type CoreConfig struct {
	MQTT struct {
		Host           string `yaml:"host" validate:"required,hostname|ip"`
		Port           string `yaml:"port" validate:"required,numeric"`
		TLS            bool   `yaml:"tls"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		BaseTopic      string `yaml:"base_topic" default:"system-bridge" validate:"required"`
		DiscoveryTopic string `yaml:"discovery_topic" default:"homeassistant" validate:"required"`
	} `yaml:"mqtt"`
	LogLevel     string `yaml:"log_level" default:"info" validate:"required,oneof=trace debug info warn error fatal panic"`
	LogLevelMqtt string `yaml:"log_level_mqtt" default:"error" validate:"required,oneof=critical error warn debug"`
}

func (c *CoreConfig) AvailabilityTopic() string {
	return path.Join(Config.MQTT.BaseTopic, HostID, "availability")
}

func Path() string {
	if env := os.Getenv("SYSTEM_BRIDGE_CONFIG"); env != "" {
		return env
	}
	return filepath.Join(utils.ExeDirectory, "config.yaml")
}

func LoadConfig() {
	// set defaults
	err := defaults.Set(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set config defaults")
	}

	// load config
	data, err := os.ReadFile(Path())
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading config")
	}

	validate := validator.New()

	// parse config
	if err = yaml.UnmarshalWithOptions(data, &Config, yaml.Validator(validate)); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing config")
	}

	log.Info().Msg("config loaded")
}
