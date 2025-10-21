package buttons

import (
	"github.com/jcwillox/system-link/utils"
	"os/exec"
	"runtime"
)
import mqtt "github.com/eclipse/paho.mqtt.golang"
import "github.com/go-co-op/gocron/v2"
import "github.com/jcwillox/system-link/entity"
import "github.com/rs/zerolog/log"

func NewLock(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("lock").
		Icon("mdi:lock").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			var command *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				command = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
			case "darwin":
				command = exec.Command("pmset", "displaysleepnow")
			case "linux":
				if utils.IsSystemd() {
					command = exec.Command("loginctl", "lock-session")
				} else {
					log.Error().Msg("lock is not supported on this linux distribution")
					return
				}
			}

			if command == nil {
				log.Error().Msg("lock is not supported on this operating system")
				return
			}

			err := command.Run()
			if err != nil {
				log.Err(err).Msg("failed to run lock command")
			}
		}).Build()
}
