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

func NewForceShutdown(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("force_shutdown").
		Icon("mdi:power").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			var command *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				command = exec.Command("shutdown", "/p")
			case "darwin":
				command = exec.Command("shutdown", "-h", "now")
			case "linux":
				if utils.IsSystemd() {
					command = exec.Command("systemctl", "poweroff", "-i")
				} else {
					command = exec.Command("shutdown", "-h", "now")
				}
			}

			if command == nil {
				log.Error().Msg("force_shutdown is not supported on this operating system")
				return
			}

			err := command.Run()
			if err != nil {
				log.Err(err).Msg("failed to run shutdown command")
			}
		}).Build()
}
