package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"syscall"
)

func NewReload(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("reload").
		Icon("mdi:restart").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Info().Msg("reloading system-link")
			if runtime.GOOS == "windows" {
				_, err := os.StartProcess(utils.ExePath, os.Args, &os.ProcAttr{
					Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
				})
				if err != nil {
					log.Err(err).Msg("failed to automatically restart system-link")
					return
				}
				os.Exit(0)
			} else {
				// exec is preferred over StartProcess as it will replace the current process,
				// but it is not available on windows
				err := syscall.Exec(utils.ExePath, os.Args, os.Environ())
				if err != nil {
					log.Err(err).Msg("failed to automatically restart system-link")
					return
				}
			}
		}).Build()
}
