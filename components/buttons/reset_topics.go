package buttons

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
	"time"
)

func NewResetTopics(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainButton).
		ID("reset_topics").
		Icon("mdi:trash-can").
		Name("Reset topics").
		OnCommand(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
			log.Info().Msg("cleaning up topics")

			messageHandler := func(client mqtt.Client, msg mqtt.Message) {
				if msg.Retained() {
					log.Info().Str("topic", msg.Topic()).Msg("cleaning up topic")
					if token := client.Publish(msg.Topic(), 0, true, ""); token.Wait() && token.Error() != nil {
						log.Err(token.Error()).Str("topic", msg.Topic()).Msg("failed to clean up topic")
						return
					}
				}
			}

			topics := []string{
				config.Config.MQTT.BaseTopic + "/+/" + config.Config.HostID + "/#",
				config.Config.MQTT.DiscoveryTopic + "/+/" + config.Config.HostID + "/#",
			}
			for _, topic := range topics {
				if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
					log.Err(token.Error()).Str("topic", topic).Msg("failed to subscribe to topic")
					return
				}
			}

			// wait for retained messages to come in
			time.Sleep(2 * time.Second)

			if token := client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
				log.Err(token.Error()).Msg("failed to unsubscribe from topics")
				return
			}

			log.Info().Msg("finished cleaning up topics, restarting...")

			// restart so we can resubscribe to correct topics
			_ = utils.RestartSelf()
		}).Build()
}
