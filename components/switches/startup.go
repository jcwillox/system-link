package switches

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
)

func createStartupEntry() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	vbsPath := filepath.Join(cwd, "system-bridge-hidden.vbs")
	err = os.WriteFile(vbsPath, config.FileSystemBridgeHiddenVBS, 0777)
	if err != nil {
		return err
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	wScriptPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "WScript.exe")
	err = key.SetStringValue("SystemBridge", `"`+wScriptPath+`" "`+vbsPath+`"`)
	if err != nil {
		return err
	}

	return nil
}

func deleteStartupEntry() error {
	// delete vbs file from startup registry
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.DeleteValue("SystemBridge")
	if err != nil {
		return err
	}

	return nil
}

func hasStartupEntry() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer key.Close()

	_, _, err = key.GetStringValue("SystemBridge")
	if err != nil {
		return false, nil
	}

	return true, nil
}

func NewStartup(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSwitch).
		ID("startup").
		Name("Run on boot").
		Icon("mdi:restart").
		EntityCategory("config").
		OnState(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			// print message
			log.Info().Bytes("payload", message.Payload()).Msg("Startup:OnState")
			// check changed and update
		}).
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Info().Bytes("payload", message.Payload()).Msg("Startup:OnCommand")
			if string(message.Payload()) == "ON" {
				err := createStartupEntry()
				if err != nil {
					log.Err(err).Msg("failed to create startup entry")
					return
				}
				_ = entity.PublishState(client, "ON")
			} else if string(message.Payload()) == "OFF" {
				err := deleteStartupEntry()
				if err != nil {
					log.Err(err).Msg("failed to delete startup entry")
					return
				}
				_ = entity.PublishState(client, "OFF")
			}
		}).
		OnSetup(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			hasEntry, err := hasStartupEntry()
			if err != nil {
				log.Err(err).Msg("failed to check for startup entry")
				return err
			}
			if hasEntry {
				return entity.PublishState(client, "ON")
			} else {
				return entity.PublishState(client, "OFF")
			}
		}).
		Build()
}
