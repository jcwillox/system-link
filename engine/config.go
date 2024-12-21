package engine

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/jcwillox/system-bridge/components/switches"
	"github.com/jcwillox/system-bridge/components/updaters"
	"github.com/rs/zerolog/log"
	"os"
)

type EntitiesConfig struct {
	Buttons       []buttons.Config        `yaml:"buttons"`
	Sensors       []sensors.Config        `yaml:"sensors"`
	BinarySensors []binary_sensors.Config `yaml:"binary_sensors"`
	Switches      []switches.Config       `yaml:"switches"`
	Updaters      []updaters.Config       `yaml:"updaters"`
}

type Config struct {
	Entities EntitiesConfig `yaml:"entities"`
}

func LoadEntitiesConfig() EntitiesConfig {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading entities config")
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing entities config")
	}

	return cfg.Entities
}
