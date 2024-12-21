package switches

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Cron    *CronConfig    `yaml:"cron,omitempty"`
	Startup *entity.Config `yaml:"startup,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Startup != nil {
		if e := NewStartup(*c.Startup); e != nil {
			entities = append(entities, e)
		}
	}

	return entities
}
