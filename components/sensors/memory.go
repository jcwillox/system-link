package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/shirou/gopsutil/v4/mem"
	"math"
)

const Gibibyte = 1024 * 1024 * 1024

func NewMemory(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("memory").
		Name("Memory").
		Icon("mdi:memory").
		StateClass("measurement").
		Unit("%").
		Precision(0).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			memory, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			used := memory.Total - memory.Available
			percent := float64(used) / float64(memory.Total) * 100
			return e.PublishState(client, math.Round(percent*10)/10)
		}).Build()
}

func NewMemoryUsed(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("memory_used").
		Name("Memory Used").
		Icon("mdi:memory").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			memory, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			used := memory.Total - memory.Available
			return e.PublishState(client, float64(used)/Gibibyte)
		}).Build()
}

func NewMemoryFree(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("memory_free").
		Name("Memory Free").
		Icon("mdi:memory").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			memory, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			return e.PublishState(client, float64(memory.Available)/Gibibyte)
		}).Build()
}
