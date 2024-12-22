package switches

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/rs/zerolog/log"
)

type CronConfig struct {
	Schedule            string `yaml:"schedule,omitempty" validate:"required,cron"`
	utils.CommandConfig `yaml:",inline"`
	entity.Config       `yaml:",inline"`
}

func NewCron(cfg CronConfig) *entity.Entity {
	if cfg.UniqueID == "" {
		log.Fatal().Msg("cron switch is missing the unique id property")
	}

	jobID := uuid.New()
	jobTask := gocron.NewTask(func() {
		log.Debug().Str("command", cfg.CommandConfig.Command).Msg("running command task")
		// run command
		err := utils.RunCommand(cfg.Command, cfg.Shell, cfg.Hidden, cfg.ShowErrors, cfg.Detached)
		if err != nil {
			log.Err(err).Str("command", cfg.CommandConfig.Command).Msg("failed to run command")
		}
	})

	updateJob := func(scheduler gocron.Scheduler) {
		log.Debug().Str("schedule", cfg.Schedule).Msg("creating cron job")
		_, err := scheduler.Update(jobID, gocron.CronJob(cfg.Schedule, true), jobTask, gocron.WithIdentifier(jobID))
		if err != nil {
			log.Err(err).Str("schedule", cfg.Schedule).Msg("failed to create cron job")
		}
	}
	removeJob := func(scheduler gocron.Scheduler) {
		log.Debug().Str("schedule", cfg.Schedule).Msg("removing cron job")
		err := scheduler.RemoveJob(jobID)
		if err != nil && !errors.Is(err, gocron.ErrJobNotFound) {
			log.Err(err).Str("schedule", cfg.Schedule).Msg("failed to remove cron job")
		}
	}

	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSwitch).
		ObjectID("cron_" + cfg.UniqueID).
		OnState(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Debug().Str("state", string(message.Payload())).Msg("cron:received state")
			if string(message.Payload()) == "ON" {
				updateJob(scheduler)
			} else if string(message.Payload()) == "OFF" {
				removeJob(scheduler)
			}
		}).
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Debug().Str("command", string(message.Payload())).Msg("cron:received command")
			if string(message.Payload()) == "ON" {
				updateJob(scheduler)
				_ = entity.PublishState(client, "ON")
			} else if string(message.Payload()) == "OFF" {
				removeJob(scheduler)
				_ = entity.PublishState(client, "OFF")
			}
		}).
		OnSetup(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			entity.OnCleanup(func() {
				removeJob(scheduler)
			})
			return nil
		}).Build()
}
