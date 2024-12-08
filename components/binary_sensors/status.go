package binary_sensors

import (
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/entity"
)

func NewStatus(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainBinarySensor).
		ID("status").
		EntityCategory("diagnostic").
		DeviceClass("connectivity").
		PayloadOn("online").
		PayloadOff("offline").
		// repurpose availability as state
		StateTopic(config.Config.AvailabilityTopic()).
		DisableAvailability().
		Build()
}
