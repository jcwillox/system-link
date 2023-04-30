package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	. "github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/engine"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

type DiscoveryConfig struct {
	Availability []struct {
		Topic string `json:"topic"`
	} `json:"availability"`
	Device DeviceConfig `json:"device"`
}

func main() {
	// ensure single instance
	// load config
	//deviceConfig := GetDeviceConfig()
	availabilityTopic := path.Join(Config.MQTT.BaseTopic, HostID, "availability")
	log.Debug().Str("topic", availabilityTopic).Msg("availability topic")

	// create mqtt client
	opts := engine.NewClientOptions()

	// set will topic
	opts.SetWill(availabilityTopic, "offline", 0, true)

	// create entities list
	entities := engine.LoadEntities()

	// create on connect handler
	//   send birth message
	//   subscribe to topics / mount entities
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Debug().Msg("mqtt: on connect")

		// publish availability
		token := client.Publish(availabilityTopic, 0, true, "online")
		if token.Wait() && token.Error() != nil {
			log.Err(token.Error()).Msg("failed publishing availability")
		} else {
			log.Debug().Msg("sent availability")
		}

		for _, entity := range entities {
			// publish config
			e := entity.GetEntity()

			//log.Debug().
			//	Str("component", e.ComponentType).
			//	Str("state_topic", e.StateTopic).
			//	Str("command_topic", e.CommandTopic).
			//	Str("unique_id", e.UniqueID).
			//	Str("name", e.Name).
			//	Msg("item")

			// marshal config
			data, err := json.Marshal(entity)
			if err != nil {
				log.Err(err).Msg("failed to marshal item")
			}

			// publish config
			token = client.Publish(e.ConfigTopic, 0, true, data)
			if token.Wait() && token.Error() != nil {
				log.Err(token.Error()).Str("name", e.Name).Msg("failed publishing config")
			} else {
				//log.Debug().Str("name", e.Name).Msg("sent config")
				//pp.Println(entity)
			}

			// subscribe to command topics
			if e.CommandTopic != "" {
				token = client.Subscribe(e.CommandTopic, 0, e.OnCommand)
				if token.Wait() && token.Error() != nil {
					log.Err(token.Error()).Str("name", e.Name).Msg("failed subscribing to command topic")
				} else {
					log.Info().Str("name", e.Name).Msg("subscribed to command topic")
				}
			}
		}
	})

	// connect to mqtt
	client := mqtt.NewClient(opts)
	log.Info().Msg("connecting to mqtt")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic().Err(token.Error()).Msg("failed connecting to broken")
	} else {
		log.Info().Msg("connected to mqtt")
	}

	// start state machine
	scheduler := gocron.NewScheduler(time.Local)

	for _, item := range entities {
		e := item.GetEntity()

		if e.StateTopic != "" {
			// TODO UpdateFrequency
			_, err := scheduler.Every(30).Seconds().Do(e.OnUpdate, client)
			if err != nil {
				log.Err(err).Str("name", e.Name).Msg("failed to schedule update")
			}
		}
	}

	scheduler.StartAsync()

	// listen for stop signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	//case <-disconnectMQTT:
	case <-sig:
	}

	// send will on graceful disconnect
	token := client.Publish(availabilityTopic, 0, true, "offline")
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Msg("failed publishing graceful will message")
	}

	log.Info().Msg("disconnecting")

	client.Disconnect(250)

	// add auto run at startup
	// add auto update
	// add filters for sensors
	//   create delta filter
	//   create throttle filter
	// add battery mode / power saving mode / change update frequency
	// create default config
	// set mqtt loggers / enable logging

	//cmd.Execute()
}
