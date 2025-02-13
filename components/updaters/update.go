package updaters

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils/update"
	"github.com/rs/zerolog/log"
	"time"
)

func NewUpdate(cfg entity.Config) *entity.Entity {
	var lastState map[string]interface{}
	return entity.NewEntity(cfg).
		Type(entity.DomainUpdate).
		ID("update").
		EntityCategory("config").
		EntityPicture("https://api.iconify.design/mdi-bridge.svg?color=%2300acff&height=96").
		PayloadInstall("install").
		OnCommand(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			updateProgress := func(progress float64) {
				if lastState == nil {
					return
				}
				lastState["update_percentage"] = progress
				data, err := json.Marshal(lastState)
				if err != nil {
					return
				}
				_ = e.PublishRawState(client, data)
			}

			err := update.Update(updateProgress)
			if err != nil {
				log.Err(err).Msg("failed to update")
			}
		}).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			log.Debug().Msg("checking for updates")

			latestVersion, err := update.GetLatestVersion()
			if err != nil {
				return err
			}

			log.Debug().Str("latest_version", latestVersion).Msg("latest version")

			lastState = map[string]interface{}{
				"installed_version": config.Version,
				"latest_version":    latestVersion,
				"title":             "System Link",
				"release_url":       config.RepoUrl + "/releases/tag/v" + latestVersion,
			}
			data, err := json.Marshal(lastState)
			if err != nil {
				return err
			}

			return e.PublishRawState(client, data)
		}).
		Interval(time.Hour).
		Build()
}
