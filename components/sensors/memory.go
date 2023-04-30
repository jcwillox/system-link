package sensors

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/mem"
)

const Gibibyte = 1024 * 1024 * 1024

type MemoryConfig = components.EntityConfig

func NewMemory(cfg MemoryConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.ObjectID = "memory"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "%"
	e.SuggestedDisplayPrecision = 1

	e.SetName("Memory")
	e.SetStateHandler(func() (interface{}, error) {
		memory, err := mem.VirtualMemory()
		if err != nil {
			return nil, err
		}
		return memory.UsedPercent, nil
	})

	e.SetDynamicOptions()
	return e
}

func NewMemoryUsed(cfg MemoryConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "memory_used"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName("Memory Used")
	e.SetStateHandler(func() (interface{}, error) {
		memory, err := mem.VirtualMemory()
		if err != nil {
			return nil, err
		}
		log.Info().Uint64("used", memory.Used).Uint64("free", memory.Free).Msg("memory")
		return float64(memory.Used) / Gibibyte, nil
	})

	e.SetDynamicOptions()
	return e
}

func NewMemoryFree(cfg MemoryConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "memory_free"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName("Memory Free")
	e.SetStateHandler(func() (interface{}, error) {
		memory, err := mem.VirtualMemory()
		if err != nil {
			return nil, err
		}
		return float64(memory.Free) / Gibibyte, nil
	})

	e.SetDynamicOptions()
	return e
}
