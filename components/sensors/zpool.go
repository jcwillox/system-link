package sensors

import (
	"encoding/json"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
)

type ZpoolConfig struct {
	Pools         *[]string `yaml:"pools"`
	entity.Config `yaml:",inline"`
}

type zpoolListOutput struct {
	Pools []zpoolInfo `json:"pools"`
}

type zpoolInfo struct {
	Name      string `json:"name"`
	Size      uint64 `json:"size"`
	Allocated uint64 `json:"allocated"`
	Free      uint64 `json:"free"`
}

func getZpoolNames(pools *[]string) ([]string, error) {
	if pools != nil {
		return *pools, nil
	}

	cmd := exec.Command("zpool", "list", "-H", "-o", "name")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var poolNames []string
	lines := string(output)
	if lines != "" {
		// Split by newlines and filter empty lines
		for _, line := range splitLines(lines) {
			if line != "" {
				poolNames = append(poolNames, line)
			}
		}
	}
	return poolNames, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func getZpoolInfo(poolName string) (*zpoolInfo, error) {
	cmd := exec.Command("zpool", "list", "-j", "--json-int", poolName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var zpoolData zpoolListOutput
	if err := json.Unmarshal(output, &zpoolData); err != nil {
		return nil, err
	}

	if len(zpoolData.Pools) == 0 {
		return nil, nil
	}

	return &zpoolData.Pools[0], nil
}

func NewZpool(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolNames, err := getZpoolNames(cfg.Pools)
	if err != nil {
		return entities
	}

	for _, poolName := range poolNames {
		pool := poolName // capture for closure
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_" + pool).
			Name(pool).
			Icon("mdi:database").
			StateClass("measurement").
			Unit("%").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				info, err := getZpoolInfo(pool)
				if err != nil {
					return err
				}
				if info == nil {
					return nil
				}
				usedPercent := float64(info.Allocated) / float64(info.Size) * 100
				return e.PublishState(client, usedPercent)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewZpoolUsed(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolNames, err := getZpoolNames(cfg.Pools)
	if err != nil {
		return entities
	}

	for _, poolName := range poolNames {
		pool := poolName // capture for closure
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_used_" + pool).
			Name(pool + " Used").
			Icon("mdi:database").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				info, err := getZpoolInfo(pool)
				if err != nil {
					return err
				}
				if info == nil {
					return nil
				}
				return e.PublishState(client, float64(info.Allocated)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewZpoolFree(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolNames, err := getZpoolNames(cfg.Pools)
	if err != nil {
		return entities
	}

	for _, poolName := range poolNames {
		pool := poolName // capture for closure
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_free_" + pool).
			Name(pool + " Free").
			Icon("mdi:database").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				info, err := getZpoolInfo(pool)
				if err != nil {
					return err
				}
				if info == nil {
					return nil
				}
				return e.PublishState(client, float64(info.Free)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}
