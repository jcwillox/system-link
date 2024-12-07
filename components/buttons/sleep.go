package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

func NewSleep(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("sleep").
		Icon("mdi:sleep").
		OnCommand(func(client mqtt.Client, message mqtt.Message) {
			err := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0").Run()
			if err != nil {
				log.Err(err).Msg("failed to run sleep command")
			}
		}).Build()
}
