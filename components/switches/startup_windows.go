//go:build windows

package switches

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/registry"
	"os"
)

func createStartupEntry() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetStringValue("SystemBridge", exePath)
}

func deleteStartupEntry() error {
	// delete vbs file from startup registry
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.DeleteValue("SystemBridge")
}

func hasStartupEntry() (bool, error) {
	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer key.Close()

	startupPath, _, err := key.GetStringValue("SystemBridge")
	if err != nil {
		return false, nil
	}

	if startupPath != exePath {
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
		DefaultStateTopic().
		OnCommand(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			if string(message.Payload()) == "ON" {
				err := createStartupEntry()
				if err != nil {
					log.Err(err).Msg("failed to create startup entry")
					return
				}
				_ = e.PublishState(client, "ON")
			} else if string(message.Payload()) == "OFF" {
				err := deleteStartupEntry()
				if err != nil {
					log.Err(err).Msg("failed to delete startup entry")
					return
				}
				_ = e.PublishState(client, "OFF")
			}
		}).
		OnSetup(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			hasEntry, err := hasStartupEntry()
			if err != nil {
				log.Err(err).Msg("failed to check for startup entry")
				return err
			}
			if hasEntry {
				return e.PublishState(client, "ON")
			} else {
				return e.PublishState(client, "OFF")
			}
		}).
		Build()
}
