package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func NewLock(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("lock").
		Icon("mdi:lock").
		OnCommand(func(client mqtt.Client, message mqtt.Message) {
			err := exec.Command("rundll32.exe", "user32.dll,LockWorkStation").Run()
			if err != nil {
				log.Err(err).Msg("failed to run lock command")
			}
		}).Build()
}
