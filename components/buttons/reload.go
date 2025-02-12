package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

func NewReload(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("reload").
		Icon("mdi:restart").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Info().Msg("reloading system-link")
			_ = utils.RestartSelf()
		}).Build()
}
