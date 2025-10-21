package buttons

import (
	"os/exec"
	"runtime"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

func NewSleep(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("sleep").
		Icon("mdi:sleep").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			var command *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				command = exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
			case "darwin":
				command = exec.Command("pmset", "sleepnow")
			case "linux":
				if utils.IsSystemd() {
					command = exec.Command("systemctl", "suspend")
				} else {
					log.Error().Msg("sleep is not supported on this linux distribution")
					return
				}
			}

			if command == nil {
				log.Error().Msg("sleep is not supported on this operating system")
				return
			}

			err := command.Run()
			if err != nil {
				log.Err(err).Msg("failed to run sleep command")
			}
		}).Build()
}
