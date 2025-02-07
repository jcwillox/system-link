package entity

import (
	"github.com/jcwillox/system-link/config"
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
	DeviceClass string   `json:"device_class,omitempty"`
	Options     []string `json:"options,omitempty"`

	StateTopic   string `json:"state_topic,omitempty"`
	ImageTopic   string `json:"image_topic,omitempty"`
	CommandTopic string `json:"command_topic,omitempty"`

	PayloadOn      string `json:"payload_on,omitempty"`
	PayloadOff     string `json:"payload_off,omitempty"`
	PayloadInstall string `json:"payload_install,omitempty"`
	PayloadLock    string `json:"payload_lock,omitempty"`
	PayloadUnlock  string `json:"payload_unlock,omitempty"`

	Availability []AvailabilityItem  `json:"availability,omitempty"`
	Device       config.DeviceConfig `json:"device,omitempty"`
	Retain       *bool               `json:"retain,omitempty"`
}

type AvailabilityItem struct {
	Topic string `json:"topic,omitempty"`
}

func (e *Entity) DiscoveryConfig() DiscoveryConfig {
	discoveryConfig := DiscoveryConfig{
		Name:                      e.Name(),
		Icon:                      e.config.Config.Icon,
		UniqueID:                  e.config.Config.UniqueID,
		EnabledByDefault:          e.config.Config.EnabledByDefault,
		JsonAttributesTopic:       e.config.Config.JsonAttributesTopic,
		ValueTemplate:             e.config.Config.ValueTemplate,
		EntityCategory:            e.config.Config.EntityCategory,
		EntityPicture:             e.config.Config.EntityPicture,
		UnitOfMeasurement:         e.config.Config.UnitOfMeasurement,
		SuggestedDisplayPrecision: e.config.Config.SuggestedDisplayPrecision,
		StateClass:                e.config.Config.StateClass,
		DeviceClass:               e.config.Config.DeviceClass,
		StateTopic:                e.StateTopic(),
		CommandTopic:              e.CommandTopic(),
		PayloadOn:                 e.config.Config.PayloadOn,
		PayloadOff:                e.config.Config.PayloadOff,
		PayloadInstall:            e.config.Config.PayloadInstall,
		PayloadLock:               e.config.Config.PayloadLock,
		PayloadUnlock:             e.config.Config.PayloadUnlock,
		Retain:                    e.config.Config.Retain,
		Options:                   e.config.Config.Options,
		Device:                    config.Device,
	}
	if e.config.componentType == DomainImage {
		discoveryConfig.ImageTopic = discoveryConfig.StateTopic
		discoveryConfig.StateTopic = ""
	}
	if e.config.Availability == nil || *e.config.Availability {
		discoveryConfig.Availability = []AvailabilityItem{{Topic: config.Config.AvailabilityTopic()}}
	}
	return discoveryConfig
}
