package buttons

import (
	"bytes"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/components"
)

type Config struct {
	Custom   *CustomConfig   `yaml:"custom,omitempty"`
	Lock     *LockConfig     `yaml:"lock,omitempty"`
	Sleep    *SleepConfig    `yaml:"sleep,omitempty"`
	Shutdown *ShutdownConfig `yaml:"shutdown,omitempty"`
}

type ButtonEntity struct {
	PressHandler func(client mqtt.Client) `json:"-"`
	*components.Entity
}

func NewButton(cfg components.EntityConfig) *ButtonEntity {
	e := &ButtonEntity{Entity: components.NewEntity(cfg)}
	e.ComponentType = "button"
	e.SetCommandHandler(func(client mqtt.Client, message mqtt.Message) {
		if bytes.Equal(message.Payload(), []byte("PRESS")) {
			e.PressHandler(client)
		}
	})
	return e
}

func (e *ButtonEntity) SetPressHandler(handler func(client mqtt.Client)) *ButtonEntity {
	e.PressHandler = handler
	return e
}
