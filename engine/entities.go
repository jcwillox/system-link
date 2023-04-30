package engine

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/jcwillox/system-bridge/components/binary_sensors"
	"github.com/jcwillox/system-bridge/components/buttons"
	"github.com/jcwillox/system-bridge/components/sensors"
	"github.com/jcwillox/system-bridge/config"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/disk"
	"time"
)

func LoadEntities() []components.EntityI {
	entities := make([]components.EntityI, 0,
		len(config.Config.Entities.Buttons)+
			len(config.Config.Entities.Sensors)+
			len(config.Config.Entities.BinarySensors),
	)

	for _, entity := range config.Config.Entities.Buttons {
		for button := range entity {
			switch button {
			case "lock":
				entities = append(entities, buttons.NewLock())
			case "sleep":
				entities = append(entities, buttons.NewSleep())
			case "shutdown":
				entities = append(entities, buttons.NewShutdown())
			default:
				log.Warn().Str("button", button).Msg("unknown button")
			}
		}
	}

	for _, entity := range config.Config.Entities.Sensors {
		entity.Uptime
		for sensor, entityConfig := range entity {
			switch sensor {
			case "cpu":
				log.Debug().Float64("update_interval", time.Duration(entityConfig.UpdateInterval).Seconds()).Msg("cpu")
				entities = append(entities, sensors.NewCPU())
			case "memory":
				entities = append(entities, sensors.NewMemory())
			case "memory_used":
				entities = append(entities, sensors.NewMemoryUsed())
			case "memory_free":
				entities = append(entities, sensors.NewMemoryFree())
			case "swap":
				entities = append(entities, sensors.NewSwap())
			case "swap_free":
				entities = append(entities, sensors.NewSwapFree())
			case "swap_used":
				entities = append(entities, sensors.NewSwapUsed())

			case "disk":
				partitions, err := disk.Partitions(false)
				if err != nil {
					return nil
				}
				for _, partition := range partitions {
					//serial, err := disk.SerialNumber(partition.Mountpoint)
					//if err != nil {
					//	log.Err(err).Str("mount", partition.Mountpoint).Msg("failed to get serial number")
					//	return nil
					//}
					//log.Info().Str("serial", serial).Str("mount", partition.Mountpoint).Msg("disk")
					label, err := disk.Label(partition.Mountpoint)
					if err != nil {
						log.Err(err).Str("mount", partition.Mountpoint).Msg("failed to get label")
						return nil
					}
					log.Info().Str("label", label).Str("mount", partition.Mountpoint).Msg("disk")
					entities = append(entities, sensors.NewDisk(label, partition.Mountpoint))
				}
			case "disk_used":
				partitions, err := disk.Partitions(false)
				if err != nil {
					return nil
				}
				for _, partition := range partitions {
					entities = append(entities, sensors.NewDiskUsed(partition.Mountpoint))
				}
			case "disk_free":
				partitions, err := disk.Partitions(false)
				if err != nil {
					return nil
				}
				for _, partition := range partitions {
					entities = append(entities, sensors.NewDiskFree(partition.Mountpoint))
				}

			case "uptime":
				entities = append(entities, sensors.NewUptime())
			default:
				log.Warn().Str("sensor", sensor).Msg("unknown sensor")
			}
		}
	}

	for _, entity := range config.Config.Entities.BinarySensors {
		for sensor := range entity {
			switch sensor {
			case "status":
				entities = append(entities, binary_sensors.NewStatus())
			default:
				log.Warn().Str("sensor", sensor).Msg("unknown binary sensor")
			}
		}
	}

	return entities
}
