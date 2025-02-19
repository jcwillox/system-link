package switches

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"time"
)

type CronEntities struct {
	Successful *entity.Config `yaml:"successful,omitempty"`
	ExitCode   *entity.Config `yaml:"exit_code,omitempty"`
	Duration   *entity.Config `yaml:"duration,omitempty"`
	Run        *entity.Config `yaml:"run,omitempty"`
}

type CronConfig struct {
	Schedule            string `yaml:"schedule,omitempty" validate:"required,cron"`
	utils.CommandConfig `yaml:",inline"`
	entity.Config       `yaml:",inline"`
	Entities            CronEntities `yaml:"entities,omitempty"`
}

func NewCron(cfg CronConfig) []*entity.Entity {
	if cfg.UniqueID == "" {
		log.Fatal().Msg("cron switch is missing the unique id property")
	}

	var (
		successEntity  *entity.Entity
		durationEntity *entity.Entity
		exitCodeEntity *entity.Entity
		entities       []*entity.Entity
	)

	if cfg.Entities.Successful != nil {
		successEntity = entity.NewEntity(*cfg.Entities.Successful).
			Type(entity.DomainBinarySensor).
			Name(cfg.Name + " Successful").
			ID(cfg.UniqueID + "_successful").
			DeviceClass("problem").
			DefaultStateTopic().Build()
		entities = append(entities, successEntity)
	}

	if cfg.Entities.Duration != nil {
		durationEntity = entity.NewEntity(*cfg.Entities.ExitCode).
			Type(entity.DomainSensor).
			Name(cfg.Name + " Duration").
			ID(cfg.UniqueID + "_duration").
			Unit("s").
			DeviceClass("duration").
			StateClass("measurement").
			Precision(2).
			DefaultStateTopic().Build()
		entities = append(entities, durationEntity)
	}

	if cfg.Entities.ExitCode != nil {
		exitCodeEntity = entity.NewEntity(*cfg.Entities.ExitCode).
			Type(entity.DomainSensor).
			Name(cfg.Name + " Exit Code").
			ID(cfg.UniqueID + "_exit_code").
			DefaultStateTopic().Build()
		entities = append(entities, exitCodeEntity)
	}

	cronEntity := entity.NewEntity(cfg.Config).
		Type(entity.DomainSwitch).
		ObjectID(cfg.UniqueID).
		ScheduleJob(gocron.CronJob(cfg.Schedule, true), func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			log.Debug().Str("command", cfg.CommandConfig.Command).Msg("running command task")
			start := time.Now()

			// run command
			res, err := utils.RunCommand(cfg.CommandConfig)

			if durationEntity != nil {
				err := durationEntity.PublishState(client, time.Since(start).Seconds())
				if err != nil {
					return err
				}
			}

			if exitCodeEntity != nil && res.Code >= 0 {
				err := exitCodeEntity.PublishState(client, res.Code)
				if err != nil {
					return err
				}
			}

			if err != nil || res.Code > 0 {
				if successEntity != nil {
					err := successEntity.PublishState(client, "ON")
					if err != nil {
						return err
					}
				}
				if err != nil {
					log.Err(err).Str("command", cfg.CommandConfig.Command).Msg("failed to run command")
					return err
				}
			} else {
				if successEntity != nil {
					err := successEntity.PublishState(client, "OFF")
					if err != nil {
						return err
					}
				}
			}

			return nil
		}).
		OnState(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
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
			_ = e.PublishState(client, string(message.Payload()))
		}).Build()

	if cfg.Entities.Run != nil {
		runEntity := entity.NewEntity(*cfg.Entities.Run).
			Type(entity.DomainButton).
			Name(cfg.Name + " Run").
			ID(cfg.UniqueID + "_run").
			OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
				_ = cronEntity.RunJob()
			}).Build()
		entities = append(entities, runEntity)
	}

	return append(entities, cronEntity)
}
