package sensors

import "github.com/jcwillox/system-bridge/components"

type Config struct {
	CPU    *CPUConfig    `yaml:"cpu,omitempty"`
	Disk   *DiskConfig   `yaml:"disk,omitempty"`
	Memory *MemoryConfig `yaml:"memory,omitempty"`
	Swap   *SwapConfig   `yaml:"swap,omitempty"`
	Uptime *UptimeConfig `yaml:"uptime,omitempty"`
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
