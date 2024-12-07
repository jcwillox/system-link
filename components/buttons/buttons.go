package buttons

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/entity"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Custom   *CustomConfig  `yaml:"custom,omitempty"`
	Lock     *entity.Config `yaml:"lock,omitempty"`
	Sleep    *entity.Config `yaml:"sleep,omitempty"`
	Shutdown *entity.Config `yaml:"shutdown,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}
