package buttons

import (
	"github.com/jcwillox/system-link/entity"
)

type Config struct {
	Custom        *CustomConfig  `yaml:"custom,omitempty"`
	Lock          *entity.Config `yaml:"lock,omitempty"`
	Sleep         *entity.Config `yaml:"sleep,omitempty"`
	Shutdown      *entity.Config `yaml:"shutdown,omitempty"`
	ForceShutdown *entity.Config `yaml:"force_shutdown,omitempty"`
	Reload        *entity.Config `yaml:"reload,omitempty"`
	Exit          *entity.Config `yaml:"exit,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Custom != nil {
		entities = append(entities, NewCustom(*c.Custom))
	}
	if c.Lock != nil {
		entities = append(entities, NewLock(*c.Lock))
	}
	if c.Shutdown != nil {
		entities = append(entities, NewShutdown(*c.Shutdown))
	}
	if c.ForceShutdown != nil {
		entities = append(entities, NewForceShutdown(*c.ForceShutdown))
	}
	if c.Sleep != nil {
		entities = append(entities, NewSleep(*c.Sleep))
	}
	if c.Reload != nil {
		entities = append(entities, NewReload(*c.Reload))
	}
	if c.Exit != nil {
		entities = append(entities, NewExit(*c.Exit))
	}

	return entities
}
