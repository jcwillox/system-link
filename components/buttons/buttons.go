package buttons

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Custom   *CustomConfig  `yaml:"custom,omitempty"`
	Lock     *entity.Config `yaml:"lock,omitempty"`
	Sleep    *entity.Config `yaml:"sleep,omitempty"`
	Shutdown *entity.Config `yaml:"shutdown,omitempty"`
	Reload   *entity.Config `yaml:"reload,omitempty"`
	Exit     *entity.Config `yaml:"exit,omitempty"`
}
