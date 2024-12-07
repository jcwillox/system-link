package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/shlex"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type CustomConfig struct {
	// The shell command to execute
	Command string `json:"command"`
	// Run command in background without opening a terminal window
	Hidden bool `json:"hidden"`
	// Whether to execute the command through the system shell
	Shell bool `json:"shell"`

	entity.Config
}

func NewCustom(cfg CustomConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainButton).
		ObjectID(cfg.UniqueID).
		OnCommand(func(client mqtt.Client, message mqtt.Message) {
			args, err := shlex.Split(cfg.Command)
			if err != nil {
				log.Err(err).Str("command", cfg.Command).Msg("failed to parse command")
			}
			err = exec.Command(args[0], args[1:]...).Run()
			if err != nil {
				log.Err(err).Str("command", cfg.Command).Msg("failed to run shutdown command")
			}
		}).Build()
}
