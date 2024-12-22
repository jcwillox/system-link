package locks

import "github.com/jcwillox/system-bridge/entity"

type Config struct {
	Custom *CustomConfig `yaml:"custom,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Custom != nil {
		entities = append(entities, NewCustom(*c.Custom))
	}

	return entities
}
