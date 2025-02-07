//go:build !windows

package switches

import (
	"github.com/jcwillox/system-link/entity"
)

func NewStartup(cfg entity.Config) *entity.Entity {
	return nil
}
