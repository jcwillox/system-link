package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-link/config"
	"github.com/jcwillox/system-link/engine"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/jcwillox/system-link/utils/update"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var version = "0.0.0"

func main() {
	// perform setup
	config.Version = version
	config.SetupLogging()
	config.LoadConfig()
	config.SetLogLevels()
	entities := engine.LoadEntities()

	// perform cleanup
	update.Cleanup()

	// ensure single instance
	lock := utils.LockAndKill()
	defer lock.Unlock()

	// setup scheduler
	scheduler, shutdown := engine.SetupScheduler()
	defer shutdown()

	// setup mqtt
	disconnect := engine.SetupMQTT(func(client mqtt.Client) {
		entity.CleanupAll(entities, client, scheduler)
		err := entity.SetupAll(entities, client, scheduler)
		if err != nil {
			log.Err(err).Msg("failed to setup entities")
		}
	})
	defer disconnect()

	// start state machine
	scheduler.Start()

	// listen for stop signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-config.ShutdownChannel:
	case <-sig:
	}
}
