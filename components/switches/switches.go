package switches

import (
	"github.com/jcwillox/go-yamltools"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cron CronConfig `yaml:"cron,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}
