//go:build linux

package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func NewShutdown(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("shutdown").
		Icon("mdi:power").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			err := exec.Command("systemctl", "poweroff").Run()
			if err != nil {
				log.Err(err).Msg("failed to run shutdown command")
			}
		}).Build()
}
