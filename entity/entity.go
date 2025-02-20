package entity

import (
	"encoding/json"
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

/* JOB UTILS */

func (e *Entity) StartJob(scheduler gocron.Scheduler) error {
	if e.config.job == nil {
		job, err := scheduler.NewJob(e.config.jobDefinition, e.config.jobTask)
		if err != nil {
			return err
		}
		e.config.job = job
	}
	return nil
}

func (e *Entity) StopJob(scheduler gocron.Scheduler) error {
	if e.config.job != nil {
		err := scheduler.RemoveJob(e.config.job.ID())
		if err != nil && !errors.Is(err, gocron.ErrJobNotFound) {
			return err
		}
		e.config.job = nil
	}
	return nil
}

func (e *Entity) UpdateJob(scheduler gocron.Scheduler, jobDefinition gocron.JobDefinition) error {
	// update job definition
	e.config.jobDefinition = jobDefinition
	// if job exists, update it with new schedule
	if e.config.job != nil {
		job, err := scheduler.Update(e.config.job.ID(), jobDefinition, e.config.jobTask)
		if err != nil {
			return err
		}
		e.config.job = job
	}
	return nil
}

func (e *Entity) RunJob() error {
	// if job exists, run it now in the scheduler
	if e.config.job != nil {
		return e.config.job.RunNow()
	}
	// otherwise invoke it directly
	e.config.jobFn()
	return nil
}

func (e *Entity) Job() gocron.Job {
	return e.config.job
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

func (e *Entity) PublishAttributes(client mqtt.Client, attributes interface{}) error {
	if e.config.attributesTopic == "" {
		log.Warn().Interface("attributes", attributes).Msg("publish failed as entity does not have an attributes topic")
		return nil
	}

	data, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	token := client.Publish(e.config.attributesTopic, 0, true, data)
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Caller().Msg("failed publishing attributes")
		return token.Error()
	}

	return nil
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
