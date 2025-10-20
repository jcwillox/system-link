//go:build linux

package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func NewLock(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("lock").
		Icon("mdi:lock").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			err := exec.Command("loginctl", "lock-session").Run()
			if err != nil {
				log.Err(err).Msg("failed to run lock command")
			}
		}).Build()
}
