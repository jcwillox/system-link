package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskConfig struct {
	Mountpoints   *[]string `yaml:"mountpoints"`
	entity.Config `yaml:",inline"`
}

func getMountpoints(mounts *[]string) []string {
	if mounts != nil {
		return *mounts
	}
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil
	}
	mountpoints := make([]string, len(partitions))
	for i, partition := range partitions {
		mountpoints[i] = partition.Mountpoint
	}
	return mountpoints
}

func NewDisk(cfg DiskConfig) []*entity.Entity {
	var entities []*entity.Entity

	for _, mount := range getMountpoints(cfg.Mountpoints) {
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("disk_" + mount).
			Name(mount).
			Icon("mdi:harddisk").
			StateClass("measurement").
			Unit("%").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				usage, err := disk.Usage(mount)
				if err != nil {
					return err
				}
				return e.PublishState(client, usage.UsedPercent)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewDiskUsed(cfg DiskConfig) []*entity.Entity {
	var entities []*entity.Entity

	for _, mount := range getMountpoints(cfg.Mountpoints) {
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("disk_used_" + mount).
			Name(mount + " Used").
			Icon("mdi:harddisk").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				usage, err := disk.Usage(mount)
				if err != nil {
					return err
				}
				return e.PublishState(client, float64(usage.Used)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewDiskFree(cfg DiskConfig) []*entity.Entity {
	var entities []*entity.Entity

	for _, mount := range getMountpoints(cfg.Mountpoints) {
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("disk_free_" + mount).
			Name(mount + " Free").
			Icon("mdi:harddisk").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				usage, err := disk.Usage(mount)
				if err != nil {
					return err
				}
				return e.PublishState(client, float64(usage.Free)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}
