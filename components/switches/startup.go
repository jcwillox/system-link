//go:build !windows

package switches

import (
	"github.com/jcwillox/system-bridge/entity"
)

func NewStartup(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainSwitch).
		ID("startup").
		Name("Run on boot").
		Icon("mdi:restart").
		EntityCategory("config").
		Build()
}
