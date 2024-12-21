package binary_sensors

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Status *entity.Config `yaml:"status,omitempty"`
}
