package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/shlex"
	"github.com/jcwillox/system-bridge/components"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type CustomConfig struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	UniqueID string `json:"unique_id"`

	// The shell command to execute
	Command string `json:"command"`
	// Run command in background without opening a terminal window
	Hidden bool `json:"hidden"`
	// Whether to execute the command through the system shell
	Shell bool `json:"shell"`

	components.EntityConfig
}

func NewCustom(cfg CustomConfig) *ButtonEntity {
	e := NewButton(cfg.EntityConfig)
	e.ObjectID = cfg.UniqueID
	e.Icon = cfg.Icon

	e.SetName(cfg.Name)
	e.SetPressHandler(func(client mqtt.Client) {
		args, err := shlex.Split(cfg.Command)
		if err != nil {
			log.Err(err).Str("command", cfg.Command).Msg("failed to parse command")
		}
		err = exec.Command(args[0], args[1:]...).Run()
		if err != nil {
			log.Err(err).Str("command", cfg.Command).Msg("failed to run shutdown command")
		}
	})

	e.SetDynamicOptions()
	return e
}
