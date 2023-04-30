package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/components"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type SleepConfig = components.EntityConfig

func NewSleep(cfg SleepConfig) *ButtonEntity {
	e := NewButton(cfg)
	e.ObjectID = "sleep"
	e.Icon = "mdi:sleep"

	e.SetName("Sleep")
	e.SetPressHandler(func(client mqtt.Client) {
		err := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0").Run()
		if err != nil {
			log.Err(err).Msg("failed to run sleep command")
		}
	})

	e.SetDynamicOptions()
	return e
}
