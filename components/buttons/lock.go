package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/components"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type LockConfig = components.EntityConfig

func NewLock(cfg LockConfig) *ButtonEntity {
	e := NewButton(cfg)
	e.ObjectID = "lock"
	e.Icon = "mdi:lock"

	e.SetName("Lock")
	e.SetPressHandler(func(client mqtt.Client) {
		err := exec.Command("rundll32.exe", "user32.dll,LockWorkStation").Run()
		if err != nil {
			log.Err(err).Msg("failed to run lock command")
		}
	})

	e.SetDynamicOptions()
	return e
}
