package sensors

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/shirou/gopsutil/v3/mem"
)

type SwapConfig = components.EntityConfig

func NewSwap(cfg SwapConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.ObjectID = "swap"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "%"
	e.SuggestedDisplayPrecision = 2

	e.SetName("Swap")
	e.SetStateHandler(func() (interface{}, error) {
		devices, err := mem.SwapDevices()
		//log.Debug().Msgf("swap devices: %v", devices)
		if err != nil {
			return "", err
		}
		if len(devices) == 0 {
			memory, err := mem.SwapMemory()
			//log.Debug().Msgf("swap memory: %v", memory)
			if err != nil {
				return "", err
			}
			return memory.UsedPercent, nil
		} else {
			total := devices[0].UsedBytes + devices[0].FreeBytes
			percent := float64(devices[0].UsedBytes) / float64(total)
			return percent, nil
		}
	})

	e.SetDynamicOptions()
	return e
}

func NewSwapUsed(cfg SwapConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "swap_used"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName("Swap Used")
	e.SetStateHandler(func() (interface{}, error) {
		devices, err := mem.SwapDevices()
		if err != nil {
			return "", err
		}
		if len(devices) == 0 {
			memory, err := mem.SwapMemory()
			if err != nil {
				return "", err
			}
			return float64(memory.Used) / Gibibyte, nil
		} else {
			return float64(devices[0].UsedBytes) / Gibibyte, nil
		}
	})

	e.SetDynamicOptions()
	return e
}

func NewSwapFree(cfg SwapConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "swap_free"
	e.Icon = "mdi:memory"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName("Swap Free")
	e.SetStateHandler(func() (interface{}, error) {
		devices, err := mem.SwapDevices()
		if err != nil {
			return "", err
		}
		if len(devices) == 0 {
			memory, err := mem.SwapMemory()
			if err != nil {
				return "", err
			}
			return float64(memory.Free) / Gibibyte, nil
		} else {
			return float64(devices[0].FreeBytes) / Gibibyte, nil
		}
	})

	e.SetDynamicOptions()
	return e
}
