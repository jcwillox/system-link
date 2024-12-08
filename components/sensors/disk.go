package sensors

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/shirou/gopsutil/v3/disk"
	"strings"
)

type DiskConfig struct {
	Mountpoint string `json:"mountpoint"`
	entity.Config
}

func NewDisk(cfg DiskConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSensor).
		ID("disk_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")).
		Name(cfg.Mountpoint).
		Icon("mdi:harddisk").
		StateClass("measurement").
		Unit("%").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			usage, err := disk.Usage(cfg.Mountpoint)
			if err != nil {
				return err
			}
			return e.PublishState(client, usage.UsedPercent)

		}).Build()
}

func NewDiskUsed(cfg DiskConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSensor).
		ID("disk_used_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")).
		Name(cfg.Mountpoint + " Used").
		Icon("mdi:harddisk").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			usage, err := disk.Usage(cfg.Mountpoint)
			if err != nil {
				return err
			}
			return e.PublishState(client, float64(usage.Used)/Gibibyte)
		}).Build()
}

func NewDiskFree(cfg DiskConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSensor).
		ID("disk_free_" + strings.ReplaceAll(cfg.Mountpoint, ":", "")).
		Name(cfg.Mountpoint + " Free").
		Icon("mdi:harddisk").
		StateClass("measurement").
		DeviceClass("data_size").
		Unit("GiB").
		Precision(1).
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			usage, err := disk.Usage(cfg.Mountpoint)
			if err != nil {
				return err
			}
			return e.PublishState(client, float64(usage.Free)/Gibibyte)
		}).Build()
}
