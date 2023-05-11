package engine

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/rs/zerolog/log"
)

type EntitiesConfig struct {
	Buttons       []buttons.Config        `yaml:"buttons"`
	Sensors       []sensors.Config        `yaml:"sensors"`
	BinarySensors []binary_sensors.Config `yaml:"binary_sensors"`
}

type Config struct {
	Entities EntitiesConfig `yaml:"entities"`
}

func LoadEntitiesConfig() EntitiesConfig {
	cfg := Config{}
	// read config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading entities config")
	}
	return cfg.Entities
}
