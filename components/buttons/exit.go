package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/entity"
)

func NewExit(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("exit").
		Icon("mdi:close-circle-outline").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			config.ShutdownChannel <- true
		}).Build()
}
