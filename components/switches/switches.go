package switches

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Cron    *CronConfig    `yaml:"cron,omitempty"`
	Startup *entity.Config `yaml:"startup,omitempty"`
}
