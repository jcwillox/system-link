package entity

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type Entity struct {
	config *BuildConfig
}

// Entity access shared root type on inherited models
func (e *Entity) Entity() *Entity {
	return e
}

/* PUBLIC ENTITY UTILS */

func (e *Entity) PublishState(client mqtt.Client, state interface{}) error {
	return e.PublishRawState(client, fmt.Sprint(state))
}

func (e *Entity) PublishRawState(client mqtt.Client, state interface{}) error {
	token := client.Publish(e.config.stateTopic, 0, true, state)
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Caller().Msg("failed publishing state")
		return token.Error()
	}
	return nil
}

/* PUBLIC GETTERS */

func (e *Entity) Name() string {
	return e.config.Config.Name
}
func (e *Entity) Icon() string {
	return e.config.Config.Icon
}
func (e *Entity) UniqueID() string {
	return e.config.Config.UniqueID
}
func (e *Entity) EnabledByDefault() *bool {
	return e.config.Config.EnabledByDefault
}
func (e *Entity) JsonAttributesTopic() string {
	return e.config.Config.JsonAttributesTopic
}
func (e *Entity) ValueTemplate() string {
	return e.config.Config.ValueTemplate
}
func (e *Entity) EntityCategory() string {
	return e.config.Config.EntityCategory
}
func (e *Entity) Unit() string {
	return e.config.Config.UnitOfMeasurement
}
func (e *Entity) Precision() int {
	return e.config.Config.SuggestedDisplayPrecision
}
func (e *Entity) StateClass() string {
	return e.config.Config.StateClass
}
func (e *Entity) DeviceClass() string {
	return e.config.Config.DeviceClass
}
func (e *Entity) StateTopic() string {
	return e.config.stateTopic
}
func (e *Entity) CommandTopic() string {
	return e.config.commandTopic
}
func (e *Entity) PayloadOn() string {
	return e.config.Config.PayloadOn
}
func (e *Entity) PayloadOff() string {
	return e.config.Config.PayloadOff
}
func (e *Entity) AvailabilityEnabled() bool {
	return !e.config.disableAvailability
}

/* INTERNAL AGGREGATE FUNCTIONS */

func SetupAll(entities []*Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
	for _, entity := range entities {
		for _, fn := range entity.config.setupFns {
			err := fn(entity, client, scheduler)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanupAll(entities []*Entity) {
	for _, entity := range entities {
		for _, fn := range entity.config.cleanupFns {
			fn()
		}
	}
}
