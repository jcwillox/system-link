package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/components"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type ShutdownConfig = components.EntityConfig

func NewShutdown(cfg ShutdownConfig) *ButtonEntity {
	e := NewButton(cfg)
	e.ObjectID = "shutdown"
	e.Icon = "mdi:power"

	e.SetName("Shutdown")
	e.SetPressHandler(func(client mqtt.Client) {
		err := exec.Command("shutdown", "/s", "/t", "0").Run()
		if err != nil {
			log.Err(err).Msg("failed to run shutdown command")
		}
	})

	e.SetDynamicOptions()
	return e
}
