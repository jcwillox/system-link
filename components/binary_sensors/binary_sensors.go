package binary_sensors

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Status *entity.Config `yaml:"status,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Status != nil {
		entities = append(entities, NewStatus(*c.Status))
	}

	return entities
}
