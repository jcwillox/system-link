package sensors

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/shirou/gopsutil/v4/cpu"
)

var ErrNoCPU = errors.New("no cpu found")

func NewCPU(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("cpu").
		Name("CPU").
		Icon("mdi:cpu-64-bit").
		StateClass("measurement").
		Unit("%").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			percent, err := cpu.Percent(0, false)
			if err != nil {
				return err
			}
			if len(percent) == 0 {
				return ErrNoCPU
			}
			return e.PublishState(client, percent[0])
		}).Build()
}
