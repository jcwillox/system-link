package updaters

import (
	"github.com/jcwillox/system-link/entity"
)

type Config struct {
	Update *entity.Config `yaml:"update,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Update != nil {
		entities = append(entities, NewUpdate(*c.Update))
	}

	return entities
}
