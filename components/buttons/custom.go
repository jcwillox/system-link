package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

type CustomConfig struct {
	utils.CommandConfig `yaml:",inline"`
	entity.Config       `yaml:",inline"`
}

func NewCustom(cfg CustomConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainButton).
		ObjectID(cfg.UniqueID).
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			_, err := utils.RunCommand(cfg.CommandConfig)
			if err != nil {
				log.Err(err).Str("command", cfg.CommandConfig.Command).Msg("failed to run command")
			}
		}).Build()
}
