package sensors

import (
	"github.com/jcwillox/go-yamltools"
	"github.com/jcwillox/system-bridge/components"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CPU *CPUConfig `yaml:"cpu,omitempty"`

	Disk     *DiskConfig `yaml:"disk,omitempty"`
	DiskUsed *DiskConfig `yaml:"disk_used,omitempty"`
	DiskFree *DiskConfig `yaml:"disk_free,omitempty"`

	Memory     *MemoryConfig `yaml:"memory,omitempty"`
	MemoryUsed *MemoryConfig `yaml:"memory_used,omitempty"`
	MemoryFree *MemoryConfig `yaml:"memory_free,omitempty"`

	Swap     *SwapConfig `yaml:"swap,omitempty"`
	SwapUsed *SwapConfig `yaml:"swap_used,omitempty"`
	SwapFree *SwapConfig `yaml:"swap_free,omitempty"`

	Uptime *UptimeConfig `yaml:"uptime,omitempty"`
}

func (c *Config) UnmarshalYAML(n *yaml.Node) error {
	n = yamltools.EnsureMapMap(n)
	type ConfigT Config
	return n.Decode((*ConfigT)(c))
}

type SensorEntity struct {
	UnitOfMeasurement         string `json:"unit_of_measurement,omitempty"`
	SuggestedDisplayPrecision int    `json:"suggested_display_precision,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`

	*components.Entity
}

func NewSensor(cfg components.EntityConfig) *SensorEntity {
	e := &SensorEntity{Entity: components.NewEntity(cfg)}
	e.ComponentType = "sensor"
	return e
}

func (e *SensorEntity) SetDynamicOptions() *SensorEntity {
	e.Entity.SetDynamicOptions()
	return e
}
