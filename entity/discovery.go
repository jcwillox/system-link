package entity

import (
	"github.com/jcwillox/system-bridge/config"
)

type DiscoveryConfig struct {
	Name             string `json:"name,omitempty"`
	Icon             string `json:"icon,omitempty"`
	UniqueID         string `json:"unique_id,omitempty"`
	EnabledByDefault *bool  `json:"enabled_by_default,omitempty"`

	JsonAttributesTopic string `json:"json_attributes_topic,omitempty"`
	ValueTemplate       string `json:"value_template,omitempty"`

	EntityCategory string `json:"entity_category,omitempty"`
	EntityPicture  string `json:"entity_picture,omitempty"`

	UnitOfMeasurement         string `json:"unit_of_measurement,omitempty"`
	SuggestedDisplayPrecision int    `json:"suggested_display_precision,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`
	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`

	StateTopic   string `json:"state_topic,omitempty"`
	ImageTopic   string `json:"image_topic,omitempty"`
	CommandTopic string `json:"command_topic,omitempty"`

	PayloadOn      string `json:"payload_on,omitempty"`
	PayloadOff     string `json:"payload_off,omitempty"`
	PayloadInstall string `json:"payload_install,omitempty"`

	Availability []AvailabilityItem  `json:"availability,omitempty"`
	Device       config.DeviceConfig `json:"device,omitempty"`
}

type AvailabilityItem struct {
	Topic string `json:"topic,omitempty"`
}

func (e *Entity) DiscoveryConfig() DiscoveryConfig {
	discoveryConfig := DiscoveryConfig{
		Name:                      e.Name(),
		Icon:                      e.Icon(),
		UniqueID:                  e.UniqueID(),
		EnabledByDefault:          e.EnabledByDefault(),
		JsonAttributesTopic:       e.JsonAttributesTopic(),
		ValueTemplate:             e.ValueTemplate(),
		EntityCategory:            e.EntityCategory(),
		EntityPicture:             e.EntityPicture(),
		UnitOfMeasurement:         e.Unit(),
		SuggestedDisplayPrecision: e.Precision(),
		StateClass:                e.StateClass(),
		DeviceClass:               e.DeviceClass(),
		StateTopic:                e.StateTopic(),
		CommandTopic:              e.CommandTopic(),
		PayloadOn:                 e.PayloadOn(),
		PayloadOff:                e.PayloadOff(),
		PayloadInstall:            e.PayloadInstall(),
		Device:                    config.Device,
	}
	if e.config.componentType == DomainImage {
		discoveryConfig.ImageTopic = discoveryConfig.StateTopic
		discoveryConfig.StateTopic = ""
	}
	if e.AvailabilityEnabled() {
		discoveryConfig.Availability = []AvailabilityItem{{Topic: config.Config.AvailabilityTopic()}}
	}
	return discoveryConfig
}
