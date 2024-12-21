package sensors

import (
	"github.com/jcwillox/system-bridge/entity"
	"github.com/shirou/gopsutil/v4/disk"
)

type Config struct {
	CPU *entity.Config `yaml:"cpu,omitempty"`

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
}

func (c *Config) LoadEntities() []*entity.Entity {
	var entities []*entity.Entity

	if c.CPU != nil {
		entities = append(entities, NewCPU(*c.CPU))
	}

	if c.Disk != nil {
		partitions, err := disk.Partitions(false)
		if err != nil {
			return nil
		}
		for _, partition := range partitions {
			c.Disk.Mountpoint = partition.Mountpoint
			entities = append(entities, NewDisk(*c.Disk))
		}
	}
	if c.DiskUsed != nil {
		partitions, err := disk.Partitions(false)
		if err != nil {
			return nil
		}
		for _, partition := range partitions {
			c.DiskUsed.Mountpoint = partition.Mountpoint
			entities = append(entities, NewDiskUsed(*c.DiskUsed))
		}
	}
	if c.DiskFree != nil {
		partitions, err := disk.Partitions(false)
		if err != nil {
			return nil
		}
		for _, partition := range partitions {
			c.DiskFree.Mountpoint = partition.Mountpoint
			entities = append(entities, NewDiskFree(*c.DiskFree))
		}
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

	return entities
}
