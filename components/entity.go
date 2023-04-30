package components

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"path"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	var s string
	err := n.Decode(&s)
	if err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

type EntityConfig struct {
	//Filters        engine.Filters `json:"filters,omitempty"`
	UpdateInterval Duration `yaml:"update_interval"`
}

type EntityI interface {
	GetEntity() *Entity
}

type Entity struct {
	// internal
	ComponentType string                      `json:"-"`
	ObjectID      string                      `json:"-"`
	ConfigTopic   string                      `json:"-"`
	OnCommand     mqtt.MessageHandler         `json:"-"`
	State         func() (interface{}, error) `json:"-"`

	// automatically handled
	Availability []AvailabilityItem  `json:"availability,omitempty"`
	Device       config.DeviceConfig `json:"device,omitempty"`

	// common attributes
	Name             string `json:"name,omitempty"`
	Icon             string `json:"icon,omitempty"`
	UniqueID         string `json:"unique_id,omitempty"`
	EnabledByDefault bool   `json:"enabled_by_default,omitempty"`

	StateTopic   string `json:"state_topic,omitempty"`
	CommandTopic string `json:"command_topic,omitempty"`

	JsonAttributesTopic string `json:"json_attributes_topic,omitempty"`
	ValueTemplate       string `json:"value_template,omitempty"`

	EntityCategory string `json:"entity_category,omitempty"`

	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`
}

type AvailabilityItem struct {
	Topic string `json:"topic,omitempty"`
}

func NewEntity(cfg EntityConfig) *Entity {
	e := &Entity{}
	e.EnabledByDefault = true
	e.Device = config.Device
	e.Availability = []AvailabilityItem{
		{Topic: path.Join(config.Config.MQTT.BaseTopic, config.HostID, "availability")},
	}
	return e
}

func (e *Entity) GetEntity() *Entity {
	return e
}

func (e *Entity) SetName(name string) *Entity {
	e.Name = config.Device.Name + " " + name
	return e
}

func (e *Entity) SetCommandHandler(handler mqtt.MessageHandler) *Entity {
	e.OnCommand = handler
	return e
}

func (e *Entity) SetStateHandler(handler func() (interface{}, error)) *Entity {
	e.State = handler
	return e
}

func (e *Entity) SetDynamicOptions() *Entity {
	switch e.ComponentType {
	case "sensor", "binary_sensor":
		e.StateTopic = path.Join(config.Config.MQTT.BaseTopic, e.ComponentType, config.HostID, e.ObjectID, "state")
	}

	switch e.ComponentType {
	case "button":
		e.CommandTopic = path.Join(config.Config.MQTT.BaseTopic, e.ComponentType, config.HostID, e.ObjectID, "set")
	}

	e.ConfigTopic = path.Join(config.Config.MQTT.DiscoveryTopic, e.ComponentType, config.HostID, e.ObjectID, "config")
	e.UniqueID = fmt.Sprintf("%s_%s_%s", config.HostID, e.ObjectID, "system-bridge")
	return e
}

func (e *Entity) OnUpdate(client mqtt.Client) {
	if e.State == nil {
		return
	}
	state, err := e.State()
	if err != nil {
		log.Err(err).Msg("failed to get state")
	}
	log.Debug().Str("name", e.Name).Interface("state", state).Msg("update")

	token := client.Publish(e.StateTopic, 0, true, fmt.Sprint(state))
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Caller().Msg("failed publishing state")
	}
}
