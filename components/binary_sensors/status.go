package binary_sensors

import "github.com/jcwillox/system-bridge/components"

type StatusConfig = components.EntityConfig

func NewStatus(cfg StatusConfig) *BinarySensorEntity {
	e := NewBinarySensor(cfg)
	e.DeviceClass = "connectivity"
	e.ObjectID = "status"

	e.PayloadOn = "online"
	e.PayloadOff = "offline"

	e.SetName("Status")
	e.SetDynamicOptions()

	// repurpose availability as state
	e.StateTopic = e.Availability[0].Topic
	e.Availability = nil

	return e
}
