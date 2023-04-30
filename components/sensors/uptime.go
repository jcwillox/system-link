package sensors

import (
	"github.com/jcwillox/system-bridge/components"
	"github.com/shirou/gopsutil/v3/host"
	"time"
)

type UptimeConfig = components.EntityConfig

func NewUptime(cfg UptimeConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.ObjectID = "uptime"
	e.DeviceClass = "timestamp"

	e.SetName("Uptime")
	e.SetStateHandler(func() (interface{}, error) {
		uptime, err := host.BootTime()
		if err != nil {
			return nil, err
		}
		return time.Unix(int64(uptime), 0).Format(time.RFC3339), nil
	})

	e.SetDynamicOptions()
	return e
}
