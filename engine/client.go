package engine

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/config"
	"github.com/rs/zerolog/log"
)

type MQTTHelpers struct {
	Client mqtt.Client
}

func (c *MQTTHelpers) Connect() {
	log.Info().Msg("connecting to mqtt")
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic().Err(token.Error()).Msg("failed connecting to broken")
	} else {
		log.Info().Msg("connected to mqtt")
	}
}

func (c *MQTTHelpers) Disconnect() {
	log.Info().Msg("disconnecting")
	// send will on graceful disconnect
	c.SendDisconnect()
	// disconnect
	c.Client.Disconnect(250)
}

func (c *MQTTHelpers) SendConnect() {
	topic := config.Config.AvailabilityTopic()
	token := c.Client.Publish(topic, 0, true, "online")
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Msg("failed publishing availability")
	} else {
		log.Debug().Msg("sent availability")
	}
}

func (c *MQTTHelpers) SendDisconnect() {
	topic := config.Config.AvailabilityTopic()
	token := c.Client.Publish(topic, 0, true, "offline")
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Msg("failed publishing graceful will message")
	}
}

func SetupMQTT(onConn mqtt.OnConnectHandler) func() {
	topic := config.Config.AvailabilityTopic()
	url := fmt.Sprintf("%s:%s", config.Config.MQTT.Host, config.Config.MQTT.Port)
	if config.Config.MQTT.TLS {
		url = "ssl://" + url
	}
	// create client options
	opts := mqtt.NewClientOptions().
		AddBroker(url).
		SetClientID("system-link-" + config.Config.HostID).
		SetUsername(config.Config.MQTT.Username).
		SetPassword(config.Config.MQTT.Password).
		SetConnectRetry(true).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false, // TODO: allow no verify / custom ssl
		})
	// set will topic
	opts.SetWill(topic, "offline", 0, true)
	// create on connect handler
	var helpers MQTTHelpers
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Debug().Msg("mqtt: on connect")
		helpers.SendConnect()
		onConn(client)
	})
	// create client and helpers
	client := mqtt.NewClient(opts)
	helpers = MQTTHelpers{client}
	helpers.Connect()
	return helpers.Disconnect
}

func SetupScheduler() (gocron.Scheduler, func()) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create scheduler")
	}
	shutdown := func() {
		err := scheduler.Shutdown()
		if err != nil {
			log.Err(err).Msg("failed to shutdown scheduler")
		}
	}
	return scheduler, shutdown
}
