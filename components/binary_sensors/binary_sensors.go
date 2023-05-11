package binary_sensors

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/components"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Status *StatusConfig `yaml:"status,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}

type BinarySensorEntity struct {
	PayloadOn  string `json:"payload_on,omitempty"`
	PayloadOff string `json:"payload_off,omitempty"`
	*components.Entity
}

func NewBinarySensor(cfg components.EntityConfig) *BinarySensorEntity {
	e := &BinarySensorEntity{Entity: components.NewEntity(cfg)}
	e.ComponentType = "binary_sensor"
	return e
}
