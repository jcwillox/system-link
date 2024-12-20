package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/mem"
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
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			memory, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			return e.PublishState(client, memory.UsedPercent)
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
			log.Info().Uint64("used", memory.Used).Uint64("free", memory.Free).Msg("memory")
			return e.PublishState(client, float64(memory.Used)/Gibibyte)
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
			return e.PublishState(client, float64(memory.Free)/Gibibyte)
		}).Build()
}
