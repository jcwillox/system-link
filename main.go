package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	. "github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/engine"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/jcwillox/system-bridge/utils/update"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var version = "0.0.0"

func main() {
	Version = version
	// setup logging & load config
	update.Cleanup()

	instanceLock, err := utils.NewInstanceLock()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create instance lock")
		return
	}

	err = instanceLock.KillLockedPid()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to kill locked pid")
		return
	}

	err = instanceLock.Lock()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to lock instance")
		return
	}

	defer instanceLock.Unlock()

	availabilityTopic := Config.AvailabilityTopic()
	log.Debug().Str("topic", availabilityTopic).Msg("availability topic")

	// create mqtt client
	opts := engine.NewClientOptions()

	// set will topic
	opts.SetWill(availabilityTopic, "offline", 0, true)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create scheduler")
	}

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

		entity.CleanupAll(entities)
		err := entity.SetupAll(entities, client, scheduler)
		if err != nil {
			log.Err(err).Msg("failed to setup entities")
		}

		//log.Debug().
		//	Str("component", e.ComponentType).
		//	Str("state_topic", e.StateTopic).
		//	Str("command_topic", e.CommandTopic).
		//	Str("unique_id", e.UniqueID).
		//	Str("name", e.Name).
		//	Msg("item")
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
	scheduler.Start()

	// listen for stop signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	//case <-disconnectMQTT:
	case <-sig:
	}

	// shutdown scheduler
	err = scheduler.Shutdown()
	if err != nil {
		log.Err(err).Msg("failed to shutdown scheduler")
	}

	// send will on graceful disconnect
	token := client.Publish(availabilityTopic, 0, true, "offline")
	if token.Wait() && token.Error() != nil {
		log.Err(token.Error()).Msg("failed publishing graceful will message")
	}

	log.Info().Msg("disconnecting")

	client.Disconnect(250)
}
