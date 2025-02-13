package switches

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
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

	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSwitch).
		ObjectID("cron_"+cfg.UniqueID).
		ScheduleJob(gocron.CronJob(cfg.Schedule, true), func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			log.Debug().Str("command", cfg.CommandConfig.Command).Msg("running command task")
			// run command
			err := utils.RunCommand(cfg.Command, cfg.Shell, cfg.Hidden, cfg.ShowErrors, cfg.Detached)
			if err != nil {
				log.Err(err).Str("command", cfg.CommandConfig.Command).Msg("failed to run command")
			}
			return nil
		}).
		OnState(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Debug().Str("state", string(message.Payload())).Msg("cron:received state")
			if string(message.Payload()) == "ON" {
				err := e.StartJob(scheduler)
				if err != nil {
					log.Err(err).Msg("failed to start cron job")
					return
				}
			} else if string(message.Payload()) == "OFF" {
				err := e.StopJob(scheduler)
				if err != nil {
					log.Err(err).Msg("failed to stop cron job")
					return
				}
			}
		}).
		OnCommand(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Debug().Str("command", string(message.Payload())).Msg("cron:received command")
			if string(message.Payload()) == "ON" {
				err := e.StartJob(scheduler)
				if err != nil {
					log.Err(err).Msg("failed to start cron job")
					return
				}
				_ = e.PublishState(client, "ON")
			} else if string(message.Payload()) == "OFF" {
				err := e.StopJob(scheduler)
				if err != nil {
					log.Err(err).Msg("failed to stop cron job")
					return
				}
				_ = e.PublishState(client, "OFF")
			}
		}).Build()
}
