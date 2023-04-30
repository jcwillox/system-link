package sensors

import (
	"errors"
	"fmt"
	"github.com/jcwillox/system-bridge/components"
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

var ErrNoCPU = errors.New("no cpu found")

type CPUConfig = components.EntityConfig

func NewCPU(cfg CPUConfig) *SensorEntity {
	e := NewSensor(cfg)
	e.StateClass = "measurement"
	e.ObjectID = "cpu"
	e.Icon = "mdi:cpu-64-bit"
	e.UnitOfMeasurement = "%"
	e.SuggestedDisplayPrecision = 1

	e.SetName("CPU")
	e.SetStateHandler(func() (interface{}, error) {
		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			return nil, err
		}
		if len(percent) == 0 {
			return nil, ErrNoCPU
		}
		return fmt.Sprint(percent[0]), nil
	})

	e.SetDynamicOptions()
	return e
}
