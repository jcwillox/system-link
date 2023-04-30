package sensors

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/shirou/gopsutil/v3/disk"
	"strings"
)

type DiskConfig struct {
	Mountpoint string `json:"mountpoint"`

	components.EntityConfig
}

func NewDisk(cfg DiskConfig) *SensorEntity {
	e := NewSensor(cfg.EntityConfig)
	e.StateClass = "measurement"
	e.ObjectID = "disk_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")
	e.Icon = "mdi:harddisk"
	e.UnitOfMeasurement = "%"
	e.SuggestedDisplayPrecision = 1

	e.SetName(cfg.Mountpoint)
	e.SetStateHandler(func() (interface{}, error) {
		usage, err := disk.Usage(cfg.Mountpoint)
		if err != nil {
			return nil, err
		}
		return usage.UsedPercent, nil
	})

	e.SetDynamicOptions()
	return e
}

func NewDiskUsed(cfg DiskConfig) *SensorEntity {
	e := NewSensor(cfg.EntityConfig)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "disk_used_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")
	e.Icon = "mdi:harddisk"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName(cfg.Mountpoint + " Used")
	e.SetStateHandler(func() (interface{}, error) {
		usage, err := disk.Usage(cfg.Mountpoint)
		if err != nil {
			return nil, err
		}
		return float64(usage.Used) / Gibibyte, nil
	})

	e.SetDynamicOptions()
	return e
}

func NewDiskFree(cfg DiskConfig) *SensorEntity {
	e := NewSensor(cfg.EntityConfig)
	e.StateClass = "measurement"
	e.DeviceClass = "data_size"
	e.ObjectID = "disk_free_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")
	e.Icon = "mdi:harddisk"
	e.UnitOfMeasurement = "GiB"
	e.SuggestedDisplayPrecision = 1

	e.SetName(cfg.Mountpoint + " Free")
	e.SetStateHandler(func() (interface{}, error) {
		usage, err := disk.Usage(cfg.Mountpoint)
		if err != nil {
			return nil, err
		}
		return float64(usage.Free) / Gibibyte, nil
	})

	e.SetDynamicOptions()
	return e
}
