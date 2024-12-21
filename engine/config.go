package engine

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/jcwillox/system-bridge/components/switches"
	"github.com/jcwillox/system-bridge/components/updaters"
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/entity"
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

func (c *EntitiesConfig) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	for _, cfg := range c.Buttons {
		entities = append(entities, cfg.LoadEntities()...)
	}
	for _, cfg := range c.Sensors {
		entities = append(entities, cfg.LoadEntities()...)
	}
	for _, cfg := range c.BinarySensors {
		entities = append(entities, cfg.LoadEntities()...)
	}
	for _, cfg := range c.Switches {
		entities = append(entities, cfg.LoadEntities()...)
	}
	for _, cfg := range c.Updaters {
		entities = append(entities, cfg.LoadEntities()...)
	}

	return entities
}

type Config struct {
	Entities EntitiesConfig `yaml:"entities"`
}

func LoadEntities() []*entity.Entity {
	data, err := os.ReadFile(config.Path())
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading entities config")
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing entities config")
	}

	return cfg.Entities.LoadEntities()
}
