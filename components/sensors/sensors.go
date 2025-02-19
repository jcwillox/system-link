package sensors

import (
	"github.com/jcwillox/system-link/entity"
)

type Config struct {
	CPU *entity.Config `yaml:"cpu,omitempty"`

	Custom *CustomConfig `yaml:"custom,omitempty"`

	Disk     *DiskConfig `yaml:"disk,omitempty"`
	DiskUsed *DiskConfig `yaml:"disk_used,omitempty"`
	DiskFree *DiskConfig `yaml:"disk_free,omitempty"`

	Memory     *entity.Config `yaml:"memory,omitempty"`
	MemoryUsed *entity.Config `yaml:"memory_used,omitempty"`
	MemoryFree *entity.Config `yaml:"memory_free,omitempty"`

	Swap     *entity.Config `yaml:"swap,omitempty"`
	SwapUsed *entity.Config `yaml:"swap_used,omitempty"`
	SwapFree *entity.Config `yaml:"swap_free,omitempty"`

	Uptime *entity.Config `yaml:"uptime,omitempty"`

	Battery      *entity.Config `yaml:"battery,omitempty"`
	BatteryState *entity.Config `yaml:"battery_state,omitempty"`
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.CPU != nil {
		entities = append(entities, NewCPU(*c.CPU))
	}

	if c.Custom != nil {
		entities = append(entities, NewCustom(*c.Custom))
	}

	if c.Disk != nil {
		entities = append(entities, NewDisk(*c.Disk)...)
	}
	if c.DiskUsed != nil {
		entities = append(entities, NewDiskUsed(*c.DiskUsed)...)
	}
	if c.DiskFree != nil {
		entities = append(entities, NewDiskFree(*c.DiskFree)...)
	}

	if c.Memory != nil {
		entities = append(entities, NewMemory(*c.Memory))
	}
	if c.MemoryUsed != nil {
		entities = append(entities, NewMemoryUsed(*c.MemoryUsed))
	}
	if c.MemoryFree != nil {
		entities = append(entities, NewMemoryFree(*c.MemoryFree))
	}

	if c.Swap != nil {
		entities = append(entities, NewSwap(*c.Swap))
	}
	if c.SwapFree != nil {
		entities = append(entities, NewSwapFree(*c.SwapFree))
	}
	if c.SwapUsed != nil {
		entities = append(entities, NewSwapUsed(*c.SwapUsed))
	}

	if c.Uptime != nil {
		entities = append(entities, NewUptime(*c.Uptime))
	}

	if c.Battery != nil {
		if e := NewBattery(*c.Battery); e != nil {
			entities = append(entities, e)
		}
	}
	if c.BatteryState != nil {
		if e := NewBatteryState(*c.BatteryState); e != nil {
			entities = append(entities, e)
		}
	}

	return entities
}
