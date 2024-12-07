package binary_sensors

import (
	"github.com/jcwillox/system-bridge/config"
	"github.com/jcwillox/system-bridge/entity"
	"path"
)

func NewStatus(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainBinarySensor).
		ID("status").
		DeviceClass("connectivity").
		PayloadOn("online").
		PayloadOff("offline").
		// repurpose availability as state
		StateTopic(path.Join(config.Config.MQTT.BaseTopic, config.HostID, "availability")).
		// todo disable availability for this entity
		Build()
}
