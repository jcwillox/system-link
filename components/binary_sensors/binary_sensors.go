package binary_sensors

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/entity"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Status *entity.Config `yaml:"status,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}
