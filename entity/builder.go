package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/filters"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path"
	"strings"
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
	DomainNumber       = Domain{"number"}
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

	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
	Step float64 `json:"step,omitempty"`

	// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
	StateClass string `json:"state_class,omitempty"`

	// https://www.home-assistant.io/integrations/sensor/#device-class
	DeviceClass string `json:"device_class,omitempty"`
	// https://www.home-assistant.io/integrations/sensor.mqtt/#options
	Options []string `json:"options,omitempty"`

	filters.Filters `yaml:",inline"`
	UpdateInterval  time.Duration `yaml:"update_interval"`

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
	cleanupFns []SetupFn

	job           gocron.Job
	jobDefinition gocron.JobDefinition
	jobTask       gocron.Task

	componentType Domain
	objectID      string

	stateTopic   string
	commandTopic string

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
		e.Config.Name = cases.Title(language.English).String(strings.ReplaceAll(id, "_", " "))
	}
	id = strings.ReplaceAll(id, ":", "")
	id = strings.ReplaceAll(id, "/", "_")
	if e.objectID == "" {
		e.objectID = id
	}
	if e.Config.UniqueID == "" {
		e.Config.UniqueID = id
	}
	return e
}

func (e *BuildConfig) ObjectID(id string) *BuildConfig {
	id = strings.ReplaceAll(id, ":", "")
	id = strings.ReplaceAll(id, "/", "_")
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

func (e *BuildConfig) Min(min float64) *BuildConfig {
	if e.Config.Min == 0 {
		e.Config.Min = min
	}
	return e
}

func (e *BuildConfig) Max(max float64) *BuildConfig {
	if e.Config.Max == 0 {
		e.Config.Max = max
	}
	return e
}

func (e *BuildConfig) Step(step float64) *BuildConfig {
	if e.Config.Step == 0 {
		e.Config.Step = step
	}
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
	if e.jobDefinition != nil {
		log.Fatal().Str("name", e.Config.Name).Msg("cannot set interval after job definition")
	}
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

func (e *BuildConfig) OnCleanup(cleanup SetupFn) *BuildConfig {
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

		e.OnCleanup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			client.Unsubscribe(e.commandTopic)
			return nil
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

		e.OnCleanup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			client.Unsubscribe(e.stateTopic)
			return nil
		})

		return nil
	})
	return e
}

func (e *BuildConfig) Schedule(handler SetupFn) *BuildConfig {
	return e.Interval(30*time.Second).ScheduleJob(gocron.DurationJob(e.UpdateInterval), handler).StartJob()
}

func (e *BuildConfig) ScheduleJob(jobDefinition gocron.JobDefinition, handler SetupFn) *BuildConfig {
	e.DefaultStateTopic().
		OnSetup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			log.Debug().Str("name", e.Config.Name).Dur("interval", e.UpdateInterval).Msg("scheduling update")

			if e.jobDefinition != nil {
				log.Warn().Str("name", e.Config.Name).Msg("entity can only have one job definition")
			}

			e.jobDefinition = jobDefinition
			e.jobTask = gocron.NewTask(func() {
				start := time.Now()

				err := handler(entity, client, scheduler)

				if config.Config.LogTiming {
					log.Info().Str("name", e.Config.Name).Stringer("duration", time.Since(start)).Msg("updated")
				}
				if err != nil {
					log.Err(err).Str("name", e.Config.Name).Msg("failed to update")
				}
			})

			return nil
		}).
		OnCleanup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			if e.job != nil {
				err := scheduler.RemoveJob(e.job.ID())
				if err != nil && !errors.Is(err, gocron.ErrJobNotFound) {
					return err
				}
				e.job = nil
			}
			e.jobDefinition = nil
			e.jobTask = nil
			return nil
		})
	return e
}

func (e *BuildConfig) StartJob() *BuildConfig {
	e.OnSetup(func(entity *Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
		if e.job == nil {
			job, err := scheduler.NewJob(
				e.jobDefinition,
				e.jobTask,
				gocron.WithStartAt(gocron.WithStartImmediately()),
				gocron.WithSingletonMode(gocron.LimitModeReschedule),
			)
			if err != nil {
				return err
			}
			e.job = job
		}
		return nil
	})
	return e
}

func (e *BuildConfig) Optimistic() *BuildConfig {
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

func (e *BuildConfig) Options(options []string) *BuildConfig {
	e.Config.Options = options
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

			log.Debug().Interface("config", discoveryConfig).Msg("discovery config")

			configTopic := path.Join(config.Config.MQTT.DiscoveryTopic, e.config.componentType.String(), config.Config.HostID, e.config.objectID, "config")
			token := client.Publish(configTopic, 0, true, data)
			if token.Wait() && token.Error() != nil {
				log.Err(token.Error()).Str("name", e.Name()).Msg("failed publishing config")
			} else {
				log.Debug().Str("name", e.Name()).Msg("sent config")
			}

			return nil
		})
}

func (e *BuildConfig) Build() *Entity {
	// todo make unique id properly unique
	// todo check runtime variables exist
	if e.Config.UniqueID != "" {
		e.Config.UniqueID = fmt.Sprintf("%s-%s-%s", "system-link", config.Config.HostID, e.Config.UniqueID)
	}
	return &Entity{config: e.EnabledByDefault()}

}
