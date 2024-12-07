package sensors

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/entity"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CPU *entity.Config `yaml:"cpu,omitempty"`

	Disk     *DiskConfig `yaml:"disk,omitempty"`
	DiskUsed *DiskConfig `yaml:"disk_used,omitempty"`
	DiskFree *DiskConfig `yaml:"disk_free,omitempty"`

	Memory     *entity.Config `yaml:"memory,omitempty"`
	MemoryUsed *entity.Config `yaml:"memory_used,omitempty"`
	MemoryFree *entity.Config `yaml:"memory_free,omitempty"`

	Swap     *entity.Config `yaml:"swap,omitempty"`
	SwapUsed *entity.Config `yaml:"swap_used,omitempty"`
	SwapFree *entity.Config `yaml:"swap_free,omitempty"`

	Uptime *entity.Config `yaml:"uptime,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}
