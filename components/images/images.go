package images

import "github.com/jcwillox/system-link/entity"

type Config struct {
	Screen *ScreenConfig `yaml:"screen,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.Screen != nil {
		entities = append(entities, NewScreen(*c.Screen)...)
	}

	return entities
}
