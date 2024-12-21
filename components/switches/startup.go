//go:build !windows

package switches

import (
	"github.com/jcwillox/system-bridge/entity"
)

func NewStartup(cfg entity.Config) *entity.Entity {
	return nil
}
