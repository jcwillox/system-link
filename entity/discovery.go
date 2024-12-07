package entity

import (
	"github.com/jcwillox/system-bridge/config"
	"path"
)

type DiscoveryConfig struct {
	Name             string `json:"name,omitempty"`
	Icon             string `json:"icon,omitempty"`
	UniqueID         string `json:"unique_id,omitempty"`
	EnabledByDefault *bool  `json:"enabled_by_default,omitempty"`

	JsonAttributesTopic string `json:"json_attributes_topic,omitempty"`
	ValueTemplate       string `json:"value_template,omitempty"`

	EntityCategory string `json:"entity_category,omitempty"`

	UnitOfMeasurement         string `json:"unit_of_measurement,omitempty"`
	SuggestedDisplayPrecision int    `json:"suggested_display_precision,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`
	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`

	StateTopic   string `json:"state_topic,omitempty"`
	CommandTopic string `json:"command_topic,omitempty"`

	Availability []AvailabilityItem  `json:"availability,omitempty"`
	Device       config.DeviceConfig `json:"device,omitempty"`
}

type AvailabilityItem struct {
	Topic string `json:"topic,omitempty"`
}

func (e *Entity) DiscoveryConfig() DiscoveryConfig {
	return DiscoveryConfig{
		Name:                      e.Name(),
		Icon:                      e.Icon(),
		UniqueID:                  e.UniqueID(),
		EnabledByDefault:          e.EnabledByDefault(),
		JsonAttributesTopic:       e.JsonAttributesTopic(),
		ValueTemplate:             e.ValueTemplate(),
		EntityCategory:            e.EntityCategory(),
		UnitOfMeasurement:         e.Unit(),
		SuggestedDisplayPrecision: e.Precision(),
		StateClass:                e.StateClass(),
		DeviceClass:               e.DeviceClass(),
		StateTopic:                e.StateTopic(),
		CommandTopic:              e.CommandTopic(),
		Availability: []AvailabilityItem{
			{Topic: path.Join(config.Config.MQTT.BaseTopic, config.HostID, "availability")},
		},
		Device: config.Device,
	}
}
