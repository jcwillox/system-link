package entity

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path"
	"time"
)

type SetupFn = func(entity *Entity, client mqtt.Client, scheduler *gocron.Scheduler) error
type MessageFn = func(entity *Entity, client mqtt.Client, scheduler *gocron.Scheduler, message mqtt.Message)

type Domain struct {
	slug string
}

func (d Domain) String() string {
	return d.slug
}

var (
	DomainBinarySensor Domain = Domain{"binary_sensor"}
	DomainButton       Domain = Domain{"button"}
	DomainSensor       Domain = Domain{"sensor"}
	DomainSwitch       Domain = Domain{"switch"}
	DomainUpdate       Domain = Domain{"update"}
)

type Config struct {
	Name             string `json:"name,omitempty"`
	Icon             string `json:"icon,omitempty"`
	UniqueID         string `json:"unique_id,omitempty"`
	EnabledByDefault *bool  `json:"enabled_by_default,omitempty"`

	JsonAttributesTopic string `json:"json_attributes_topic,omitempty"`
	ValueTemplate       string `json:"value_template,omitempty"`

	// `diagnostic`, `config`
	EntityCategory string `json:"entity_category,omitempty"`

	UnitOfMeasurement         string `json:"unit_of_measurement,omitempty"`
	SuggestedDisplayPrecision int    `json:"suggested_display_precision,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`

	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`

	//Filters        engine.Filters `json:"filters,omitempty"`
	UpdateInterval utils.Duration `yaml:"update_interval"`

	PayloadOn  string `json:"payload_on,omitempty"`
	PayloadOff string `json:"payload_off,omitempty"`
}

type BuildConfig struct {
	setupFns   []SetupFn
	cleanupFns []func()

	componentType Domain
	objectID      string

	stateTopic   string
	commandTopic string

	disableAvailability bool

	Config
}

func (e *BuildConfig) DefaultStateTopic() *BuildConfig {
	if e.stateTopic == "" {
		if e.componentType == (Domain{}) || e.objectID == "" {
			log.Fatal().Str("name", e.Config.Name).Stringer("type", e.componentType).
				Str("object", e.objectID).Msg("component type and object id must be set for state topic")
		}

		e.stateTopic = path.Join(config.Config.MQTT.BaseTopic, e.componentType.String(), config.HostID, e.objectID, "state")
	}
	return e
}

func (e *BuildConfig) DefaultCommandTopic() *BuildConfig {
	if e.commandTopic == "" {
		if e.componentType == (Domain{}) || e.objectID == "" {
			log.Fatal().Str("name", e.Config.Name).Stringer("type", e.componentType).
				Str("object", e.objectID).Msg("component type and object id must be set for command topic")
		}

		e.commandTopic = path.Join(config.Config.MQTT.BaseTopic, e.componentType.String(), config.HostID, e.objectID, "set")
	}
	return e
}

func (e *BuildConfig) Type(component Domain) *BuildConfig {
	e.componentType = component
	return e
}

func (e *BuildConfig) ID(id string) *BuildConfig {
	if e.Config.Name == "" {
		e.Config.Name = cases.Title(language.English).String(id)
	}
	if e.objectID == "" {
		e.objectID = id
	}
	if e.Config.UniqueID == "" {
		e.Config.UniqueID = fmt.Sprintf("%s_%s_%s", config.HostID, id, "system-bridge")
	}
	return e
}

func (e *BuildConfig) ObjectID(id string) *BuildConfig {
	e.objectID = id
	return e
}

func (e *BuildConfig) Name(name string) *BuildConfig {
	e.Config.Name = name
	return e
}

func (e *BuildConfig) Icon(icon string) *BuildConfig {
	e.Config.Icon = icon
	return e
}

func (e *BuildConfig) Unit(unit string) *BuildConfig {
	e.Config.UnitOfMeasurement = unit
	return e
}

func (e *BuildConfig) Precision(precision int) *BuildConfig {
	e.Config.SuggestedDisplayPrecision = precision
	return e
}

func (e *BuildConfig) StateClass(class string) *BuildConfig {
	e.Config.StateClass = class
	return e
}

func (e *BuildConfig) DeviceClass(class string) *BuildConfig {
	e.Config.DeviceClass = class
	return e
}

func (e *BuildConfig) StateTopic(topic string) *BuildConfig {
	e.stateTopic = topic
	return e
}

func (e *BuildConfig) CommandTopic(topic string) *BuildConfig {
	e.commandTopic = topic
	return e
}

func (e *BuildConfig) EntityCategory(category string) *BuildConfig {
	e.Config.EntityCategory = category
	return e
}

func (e *BuildConfig) EnabledByDefault() *BuildConfig {
	if e.Config.EnabledByDefault == nil {
		e.Config.EnabledByDefault = new(bool)
		*e.Config.EnabledByDefault = true
	}
	return e
}

func (e *BuildConfig) DisabledByDefault() *BuildConfig {
	if e.Config.EnabledByDefault == nil {
		e.Config.EnabledByDefault = new(bool)
		*e.Config.EnabledByDefault = false
	}
	return e
}

func (e *BuildConfig) EnableAvailability() *BuildConfig {
	e.disableAvailability = false
	return e
}

func (e *BuildConfig) DisableAvailability() *BuildConfig {
	e.disableAvailability = true
	return e
}

func (e *BuildConfig) Interval(interval time.Duration) *BuildConfig {
	if e.Config.UpdateInterval == 0 {
		e.Config.UpdateInterval = utils.Duration(interval)
	}
	return e
}

func (e *BuildConfig) PayloadOn(payload string) *BuildConfig {
	e.Config.PayloadOn = payload
	return e
}
func (e *BuildConfig) PayloadOff(payload string) *BuildConfig {
	e.Config.PayloadOff = payload
	return e
}

/* EVENT HANDLERS */

func (e *BuildConfig) OnSetup(setupFn SetupFn) *BuildConfig {
	e.setupFns = append(e.setupFns, setupFn)
	return e
}

func (e *BuildConfig) OnCleanup(cleanup func()) *BuildConfig {
	e.cleanupFns = append(e.cleanupFns, cleanup)
	return e
}

func (e *BuildConfig) OnCommand(handler mqtt.MessageHandler) *BuildConfig {
	e.DefaultCommandTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
		token := client.Subscribe(e.commandTopic, 0, handler)
		if token.Wait() && token.Error() != nil {
			log.Err(token.Error()).Str("name", e.Config.Name).Msg("failed subscribing to command topic")
		} else {
			log.Info().Str("name", e.Config.Name).Msg("subscribed to command topic")
		}

		e.OnCleanup(func() {
			client.Unsubscribe(e.commandTopic)
		})

		return nil
	})
	return e
}

func (e *BuildConfig) OnState(handler mqtt.MessageHandler) *BuildConfig {
	e.DefaultStateTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
		token := client.Subscribe(e.stateTopic, 0, handler)
		if token.Wait() && token.Error() != nil {
			log.Err(token.Error()).Str("name", e.Config.Name).Msg("failed subscribing to state topic")
		} else {
			log.Info().Str("name", e.Config.Name).Msg("subscribed to state topic")
		}

		e.OnCleanup(func() {
			client.Unsubscribe(e.stateTopic)
		})

		return nil
	})
	return e
}

func (e *BuildConfig) Schedule(handler SetupFn) *BuildConfig {
	e.DefaultStateTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
		log.Info().Str("name", e.Config.Name).Dur("interval", time.Duration(e.UpdateInterval)).Msg("scheduling update")
		job, err := scheduler.Every(time.Duration(e.UpdateInterval)).Do(func() {
			err := handler(entity, client, scheduler)
			if err != nil {
				log.Err(err).Str("name", e.Config.Name).Msg("failed to update")
			}
		})

		if err != nil {
			log.Err(err).Str("name", e.Config.Name).Msg("failed to schedule update")
		}

		e.OnCleanup(func() {
			scheduler.RemoveByReference(job)
		})

		return nil
	})
	return e
}

func NewEntity(cfg Config) *BuildConfig {
	return (&BuildConfig{Config: cfg}).
		OnSetup(func(e *Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
			discoveryConfig := e.DiscoveryConfig()
			data, err := json.Marshal(discoveryConfig)
			if err != nil {
				log.Err(err).Msg("failed to marshal item")
			}

			log.Info().Interface("config", discoveryConfig).Msg("discovery config")

			configTopic := path.Join(config.Config.MQTT.DiscoveryTopic, e.config.componentType.String(), config.HostID, e.config.objectID, "config")
			token := client.Publish(configTopic, 0, true, data)
			if token.Wait() && token.Error() != nil {
				log.Err(token.Error()).Str("name", e.Name()).Msg("failed publishing config")
			} else {
				log.Debug().Str("name", e.Name()).Msg("sent config")
				//pp.Println(entity)
			}

			return nil
		})
}

func (e *BuildConfig) Build() *Entity {
	// todo make unique id properly unique
	// todo check runtime variables exist
	return &Entity{config: e.EnabledByDefault().Interval(30 * time.Second)}

}
