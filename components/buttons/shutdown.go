package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func NewShutdown(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("shutdown").
		Icon("mdi:power").
		OnCommand(func(client mqtt.Client, message mqtt.Message) {
			err := exec.Command("shutdown", "/s", "/t", "0").Run()
			if err != nil {
				log.Err(err).Msg("failed to run shutdown command")
			}
		}).Build()
}
