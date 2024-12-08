package switches

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/entity"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cron    *CronConfig    `yaml:"cron,omitempty"`
	Startup *entity.Config `yaml:"startup,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}
