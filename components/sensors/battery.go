package sensors

import (
	"github.com/distatus/battery"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
)

func NewBattery(cfg entity.Config) *entity.Entity {
	batteries, err := battery.GetAll()
	if err != nil {
		log.Err(err).Msg("failed to get battery info")
	}
	if len(batteries) == 0 {
		return nil
	}
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("battery").
		DeviceClass("battery").
		StateClass("measurement").
		Unit("%").
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			batteries, err := battery.GetAll()
			if err != nil {
				return err
			}
			totalCurrent := 0.0
			totalFull := 0.0
			for _, b := range batteries {
				totalCurrent += b.Current
				totalFull += b.Full
			}
			if totalFull > 0 {
				return e.PublishState(client, (totalCurrent/totalFull)*100)
			}
			return nil
		}).Build()
}

func allBattery(batteries []*battery.Battery, status battery.AgnosticState) bool {
	for _, b := range batteries {
		if b.State.Raw != status {
			return false
		}
	}
	return true
}

func anyBattery(batteries []*battery.Battery, status battery.AgnosticState) bool {
	for _, b := range batteries {
		if b.State.Raw == status {
			return true
		}
	}
	return false
}

func NewBatteryState(cfg entity.Config) *entity.Entity {
	batteries, err := battery.GetAll()
	if err != nil {
		log.Err(err).Msg("failed to get battery info")
	}
	if len(batteries) == 0 {
		return nil
	}
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("battery_state").
		Icon("mdi:battery").
		DeviceClass("enum").
		Options([]string{"charging", "discharging", "full", "not_charging"}).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			batteries, err := battery.GetAll()
			if err != nil {
				return err
			}
			if allBattery(batteries, battery.Full) {
				return e.PublishState(client, "full")
			}
			if anyBattery(batteries, battery.Discharging) {
				return e.PublishState(client, "discharging")
			}
			if anyBattery(batteries, battery.Charging) {
				return e.PublishState(client, "charging")
			}
			return e.PublishState(client, "not_charging")
		}).Build()
}
