package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/shirou/gopsutil/v3/mem"
)

func NewSwap(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("swap").
		Name("Swap").
		Icon("mdi:memory").
		StateClass("measurement").
		Unit("%").
		Precision(2).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
			devices, err := mem.SwapDevices()
			//log.Debug().Msgf("swap devices: %v", devices)
			if err != nil {
				return err
			}
			if len(devices) == 0 {
				memory, err := mem.SwapMemory()
				//log.Debug().Msgf("swap memory: %v", memory)
				if err != nil {
					return err
				}
				return e.PublishState(client, memory.UsedPercent)
			} else {
				total := devices[0].UsedBytes + devices[0].FreeBytes
				percent := float64(devices[0].UsedBytes) / float64(total)
				return e.PublishState(client, percent)
			}
		}).Build()
}

func NewSwapUsed(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("swap_used").
		Name("Swap Used").
		Icon("mdi:memory").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
			devices, err := mem.SwapDevices()
			if err != nil {
				return err
			}
			if len(devices) == 0 {
				memory, err := mem.SwapMemory()
				if err != nil {
					return err
				}
				return e.PublishState(client, float64(memory.Used)/Gibibyte)
			} else {
				return e.PublishState(client, float64(devices[0].UsedBytes)/Gibibyte)
			}
		}).Build()
}

func NewSwapFree(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSensor).
		ID("swap_free").
		Name("Swap Free").
		Icon("mdi:memory").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler *gocron.Scheduler) error {
			devices, err := mem.SwapDevices()
			if err != nil {
				return err
			}
			if len(devices) == 0 {
				memory, err := mem.SwapMemory()
				if err != nil {
					return err
				}
				return e.PublishState(client, float64(memory.Free)/Gibibyte)
			} else {
				return e.PublishState(client, float64(devices[0].FreeBytes)/Gibibyte)
			}
		}).Build()
}
