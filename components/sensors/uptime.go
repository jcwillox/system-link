package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/shirou/gopsutil/v3/host"
	"time"
)

func NewUptime(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("uptime").
		EntityCategory("diagnostic").
		DeviceClass("timestamp").
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			uptime, err := host.BootTime()
			if err != nil {
				return err
			}
			return e.PublishState(client, time.Unix(int64(uptime), 0).Format(time.RFC3339))
		}).Build()
}
