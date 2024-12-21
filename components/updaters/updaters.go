package updaters

import (
	"github.com/jcwillox/system-bridge/entity"
)

type Config struct {
	Update *entity.Config `yaml:"update,omitempty"`
}
