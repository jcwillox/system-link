package engine

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/shirou/gopsutil/v3/disk"
)

func LoadEntities() []components.EntityI {
	cfg := LoadEntitiesConfig()

	entities := make([]components.EntityI, 0,
		len(cfg.Buttons)+
			len(cfg.Sensors)+
			len(cfg.BinarySensors),
	)

	for _, entity := range cfg.Buttons {
		if entity.Custom != nil {
			entities = append(entities, buttons.NewCustom(*entity.Custom))
		}
		if entity.Lock != nil {
			entities = append(entities, buttons.NewLock(*entity.Lock))
		}
		if entity.Shutdown != nil {
			entities = append(entities, buttons.NewShutdown(*entity.Shutdown))
		}
		if entity.Sleep != nil {
			entities = append(entities, buttons.NewSleep(*entity.Sleep))
		}
	}

	for _, entity := range cfg.Sensors {
		if entity.CPU != nil {
			entities = append(entities, sensors.NewCPU(*entity.CPU))
		}

		if entity.Disk != nil {
			partitions, err := disk.Partitions(false)
			if err != nil {
				return nil
			}
			for _, partition := range partitions {
				entity.Disk.Mountpoint = partition.Mountpoint
				entities = append(entities, sensors.NewDisk(*entity.Disk))
			}
		}
		if entity.DiskUsed != nil {
			partitions, err := disk.Partitions(false)
			if err != nil {
				return nil
			}
			for _, partition := range partitions {
				entity.DiskUsed.Mountpoint = partition.Mountpoint
				entities = append(entities, sensors.NewDiskUsed(*entity.DiskUsed))
			}
		}
		if entity.DiskFree != nil {
			partitions, err := disk.Partitions(false)
			if err != nil {
				return nil
			}
			for _, partition := range partitions {
				entity.DiskFree.Mountpoint = partition.Mountpoint
				entities = append(entities, sensors.NewDiskFree(*entity.DiskFree))
			}
		}

		if entity.Memory != nil {
			entities = append(entities, sensors.NewMemory(*entity.Memory))
		}
		if entity.MemoryUsed != nil {
			entities = append(entities, sensors.NewMemoryUsed(*entity.MemoryUsed))
		}
		if entity.MemoryFree != nil {
			entities = append(entities, sensors.NewMemoryFree(*entity.MemoryFree))
		}

		if entity.Swap != nil {
			entities = append(entities, sensors.NewSwap(*entity.Swap))
		}
		if entity.SwapFree != nil {
			entities = append(entities, sensors.NewSwapFree(*entity.SwapFree))
		}
		if entity.SwapUsed != nil {
			entities = append(entities, sensors.NewSwapUsed(*entity.SwapUsed))
		}

		if entity.Uptime != nil {
			entities = append(entities, sensors.NewUptime(*entity.Uptime))
		}
	}

	for _, entity := range cfg.BinarySensors {
		if entity.Status != nil {
			entities = append(entities, binary_sensors.NewStatus(*entity.Status))
		}
	}

	return entities
}
