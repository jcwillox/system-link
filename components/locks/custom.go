package locks

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/rs/zerolog/log"
)

type CustomConfig struct {
	Lock          utils.CommandConfig `yaml:"lock"`
	Unlock        utils.CommandConfig `yaml:"unlock"`
	entity.Config `yaml:",inline"`
}

func NewCustom(cfg CustomConfig) *entity.Entity {
	fmt.Println(cfg)
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainLock).
		ObjectID(cfg.UniqueID).
		OnState(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			if string(message.Payload()) == "LOCKED" {
				err := utils.RunCommand(cfg.Lock.Command, cfg.Lock.Shell, cfg.Lock.Hidden, cfg.Lock.ShowErrors, cfg.Lock.Detached)
				if err != nil {
					log.Err(err).Str("command", cfg.Lock.Command).Msg("failed to run command")
				}
			} else if string(message.Payload()) == "UNLOCKED" {
				err := utils.RunCommand(cfg.Unlock.Command, cfg.Unlock.Shell, cfg.Unlock.Hidden, cfg.Unlock.ShowErrors, cfg.Unlock.Detached)
				if err != nil {
					log.Err(err).Str("command", cfg.Unlock.Command).Msg("failed to run command")
				}
			}
		}).
		OnCommand(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			if string(message.Payload()) == "LOCK" {
				_ = e.PublishState(client, "LOCKED")
			} else if string(message.Payload()) == "UNLOCK" {
				_ = e.PublishState(client, "UNLOCKED")
			}
		}).Build()
}
