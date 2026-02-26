package engine

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jcwillox/system-link/components/binary_sensors"
	"github.com/jcwillox/system-link/components/buttons"
	"github.com/jcwillox/system-link/components/images"
	"github.com/jcwillox/system-link/components/locks"
	"github.com/jcwillox/system-link/components/sensors"
	"github.com/jcwillox/system-link/components/switches"
	"github.com/jcwillox/system-link/components/updaters"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

type EntitiesConfig struct {
	Buttons       []buttons.Config        `yaml:"buttons"`
	Sensors       []sensors.Config        `yaml:"sensors"`
	BinarySensors []binary_sensors.Config `yaml:"binary_sensors"`
	Switches      []switches.Config       `yaml:"switches"`
	Updaters      []updaters.Config       `yaml:"updaters"`
	Locks         []locks.Config          `yaml:"locks"`
	Images        []images.Config         `yaml:"images"`
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
	for _, cfg := range c.Locks {
		entities = append(entities, cfg.LoadEntities()...)
	}
	for _, cfg := range c.Images {
		entities = append(entities, cfg.LoadEntities()...)
	}

	return entities
}

func LoadEntities() []*entity.Entity {
	data, err := os.ReadFile(utils.ConfigPath())
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error reading entities config")
	}

	cfg := EntitiesConfig{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		log.Fatal().Msg("fatal error parsing entities config")
	}

	return cfg.LoadEntities()
}
