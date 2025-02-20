package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	RepoUrl         = "https://github.com/jcwillox/system-link"
	Version         string
	Config          CoreConfig
	ShutdownChannel = make(chan bool)
	Path            = sync.OnceValue(func() string {
		// attempt to load from exe directory first
		configPath := filepath.Join(utils.ExeDirectory, "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		// try load from env path
		if env := os.Getenv("SYSTEM_LINK_CONFIG"); env != "" {
			return env
		}
		// fallback to exe directory
		return configPath
	})
)

type CoreConfig struct {
	HostID     string `yaml:"host_id"`
	DeviceName string `yaml:"device_name"`
	MQTT       struct {
		Host           string `yaml:"host" validate:"required,hostname|ip"`
		Port           string `yaml:"port" validate:"required,numeric"`
		TLS            bool   `yaml:"tls"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		BaseTopic      string `yaml:"base_topic" default:"system-link" validate:"required"`
		DiscoveryTopic string `yaml:"discovery_topic" default:"homeassistant" validate:"required"`
	} `yaml:"mqtt"`
	LogLevel     string `yaml:"log_level" default:"info" validate:"required,oneof=trace debug info warn error fatal panic"`
	LogLevelMqtt string `yaml:"log_level_mqtt" default:"error" validate:"required,oneof=critical error warn debug"`
	LogTiming    bool   `yaml:"log_timing" default:"false"`
}

func (c *CoreConfig) AvailabilityTopic() string {
	return path.Join(Config.MQTT.BaseTopic, Config.HostID, "availability")
}

func LogsPath() string {
	name := strings.TrimSuffix(utils.ExeName, ".exe") + ".log"
	if env := os.Getenv("SYSTEM_LINK_LOGS_DIR"); env != "" {
		return filepath.Join(env, name)
	}
	return filepath.Join(utils.ExeDirectory, name)
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
		log.Fatal().Err(err).Str("path", Path()).Msg("fatal error reading config")
	}

	validate := validator.New()

	// parse config
	if err = yaml.UnmarshalWithOptions(data, &Config, yaml.Validator(validate)); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing config")
	}

	// load device config
	loadDeviceConfig()

	log.Info().Msg("config loaded")
}
