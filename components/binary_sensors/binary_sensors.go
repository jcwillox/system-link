package binary_sensors

import "github.com/jcwillox/system-bridge/components"

type Config struct {
	Status *StatusConfig `yaml:"status,omitempty"`
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
