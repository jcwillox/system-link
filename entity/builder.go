package entity

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path"
	"time"
)

type SetupFn = func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error
type MessageFn = func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message)

type Domain struct {
	slug string
}

func (d Domain) String() string {
	return d.slug
}

var (
	DomainBinarySensor = Domain{"binary_sensor"}
	DomainButton       = Domain{"button"}
	DomainSensor       = Domain{"sensor"}
	DomainSwitch       = Domain{"switch"}
	DomainUpdate       = Domain{"update"}
	DomainLock         = Domain{"lock"}
	DomainImage        = Domain{"image"}
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
	EntityPicture  string `json:"entity_picture,omitempty"`

	UnitOfMeasurement         string `json:"unit_of_measurement,omitempty"`
	SuggestedDisplayPrecision int    `json:"suggested_display_precision,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`

	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`

	//Filters        engine.Filters `json:"filters,omitempty"`
	UpdateInterval time.Duration `yaml:"update_interval"`

	PayloadOn      string `json:"payload_on,omitempty"`
	PayloadOff     string `json:"payload_off,omitempty"`
	PayloadInstall string `json:"payload_install,omitempty"`
	PayloadLock    string `json:"payload_lock,omitempty"`
	PayloadUnlock  string `json:"payload_unlock,omitempty"`

	// set too false to disable availability
	Availability *bool `json:"availability,omitempty"`
	Retain       *bool `json:"retain,omitempty"`
}

type BuildConfig struct {
	setupFns   []SetupFn
	cleanupFns []func()

	componentType Domain
	objectID      string

	stateTopic   string
	commandTopic string

	runScheduleAtStart bool

	Config
}

func (e *BuildConfig) DefaultStateTopic() *BuildConfig {
	if e.stateTopic == "" {
		if e.componentType == (Domain{}) || e.objectID == "" {
			log.Fatal().Str("name", e.Config.Name).Stringer("type", e.componentType).
				Str("object", e.objectID).Msg("component type and object id must be set for state topic")
		}

		e.stateTopic = path.Join(config.Config.MQTT.BaseTopic, e.componentType.String(), config.Config.HostID, e.objectID, "state")
	}
	return e
}

func (e *BuildConfig) DefaultCommandTopic() *BuildConfig {
	if e.commandTopic == "" {
		if e.componentType == (Domain{}) || e.objectID == "" {
			log.Fatal().Str("name", e.Config.Name).Stringer("type", e.componentType).
				Str("object", e.objectID).Msg("component type and object id must be set for command topic")
		}

		e.commandTopic = path.Join(config.Config.MQTT.BaseTopic, e.componentType.String(), config.Config.HostID, e.objectID, "set")
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
		e.Config.UniqueID = fmt.Sprintf("%s_%s_%s", config.Config.HostID, id, "system-bridge")
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

func (e *BuildConfig) EntityPicture(picture string) *BuildConfig {
	e.Config.EntityPicture = picture
	return e
}

func (e *BuildConfig) PayloadInstall(payload string) *BuildConfig {
	e.Config.PayloadInstall = payload
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
	if e.Config.Availability == nil {
		e.Config.Availability = new(bool)
		*e.Config.Availability = true
	}
	return e
}

func (e *BuildConfig) DisableAvailability() *BuildConfig {
	if e.Config.Availability == nil {
		e.Config.Availability = new(bool)
		*e.Config.Availability = false
	}
	return e
}

func (e *BuildConfig) Interval(interval time.Duration) *BuildConfig {
	if e.Config.UpdateInterval == 0 {
		e.Config.UpdateInterval = interval
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

func (e *BuildConfig) PayloadLock(payload string) *BuildConfig {
	e.Config.PayloadLock = payload
	return e
}

func (e *BuildConfig) PayloadUnlock(payload string) *BuildConfig {
	e.Config.PayloadUnlock = payload
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

func (e *BuildConfig) OnCommand(handler MessageFn) *BuildConfig {
	e.DefaultCommandTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
		token := client.Subscribe(e.commandTopic, 0, func(client mqtt.Client, message mqtt.Message) {
			handler(entity, client, scheduler, message)
		})
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

func (e *BuildConfig) OnState(handler MessageFn) *BuildConfig {
	e.DefaultStateTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
		token := client.Subscribe(e.stateTopic, 0, func(client mqtt.Client, message mqtt.Message) {
			handler(entity, client, scheduler, message)
		})
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
	e.DefaultStateTopic().OnSetup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
		log.Info().Str("name", e.Config.Name).Dur("interval", time.Duration(e.UpdateInterval)).Msg("scheduling update")

		var options []gocron.JobOption

		if e.runScheduleAtStart {
			options = append(options, gocron.WithStartAt(gocron.WithStartImmediately()))
		}

		job, err := scheduler.NewJob(gocron.DurationJob(time.Duration(e.UpdateInterval)), gocron.NewTask(func() {
			err := handler(entity, client, scheduler)
			if err != nil {
				log.Err(err).Str("name", e.Config.Name).Msg("failed to update")
			}
		}), options...)

		if err != nil {
			log.Err(err).Str("name", e.Config.Name).Msg("failed to schedule update")
		}

		e.OnCleanup(func() {
			err := scheduler.RemoveJob(job.ID())
			if err != nil {
				log.Err(err).Str("name", e.Config.Name).Msg("failed to remove job")
			}
		})

		return nil
	})
	return e
}

func (e *BuildConfig) RunAtStart() *BuildConfig {
	e.runScheduleAtStart = true
	return e
}

func (e *BuildConfig) Optimistic() *BuildConfig {
	//e.DefaultStateTopic()
	e.commandTopic = e.stateTopic
	return e
}

func (e *BuildConfig) Retain() *BuildConfig {
	if e.Config.Retain == nil {
		e.Config.Retain = new(bool)
		*e.Config.Retain = true
	}
	return e
}

func NewEntity(cfg Config) *BuildConfig {
	return (&BuildConfig{Config: cfg}).
		OnSetup(func(e *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			discoveryConfig := e.DiscoveryConfig()
			data, err := json.Marshal(discoveryConfig)
			if err != nil {
				log.Err(err).Msg("failed to marshal item")
			}

			log.Info().Interface("config", discoveryConfig).Msg("discovery config")

			configTopic := path.Join(config.Config.MQTT.DiscoveryTopic, e.config.componentType.String(), config.Config.HostID, e.config.objectID, "config")
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
