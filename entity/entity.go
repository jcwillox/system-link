package entity

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/filters"
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
	var err error
	if len(e.config.Filters.Filters) > 0 {
		state, err = e.config.Filter(state)
		if err != nil {
			if errors.Is(err, filters.SkipSendErr) {
				return nil
			}
			return err
		}
	}
	return e.PublishRawState(client, fmt.Sprint(state))
}

func (e *Entity) PublishRawState(client mqtt.Client, state interface{}) error {
	if e.config.stateTopic == "" {
		log.Warn().Interface("state", state).Msg("publish failed as entity does not have a state topic")
		return nil
	}
	token := client.Publish(e.config.stateTopic, 0, true, state)
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Caller().Msg("failed publishing state")
		return token.Error()
	}
	return nil
}

func (e *Entity) OnCleanup(fn SetupFn) *Entity {
	e.config.OnCleanup(fn)
	return e
}

/* PUBLIC GETTERS */

func (e *Entity) Name() string {
	return e.config.Config.Name
}
func (e *Entity) StateTopic() string {
	return e.config.stateTopic
}
func (e *Entity) CommandTopic() string {
	return e.config.commandTopic
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

func CleanupAll(entities []*Entity, client mqtt.Client, scheduler gocron.Scheduler) {
	for _, entity := range entities {
		for _, fn := range entity.config.cleanupFns {
			err := fn(entity, client, scheduler)
			if err != nil {
				log.Err(err).Msg("failed to cleanup entity")
			}
		}
	}
}
