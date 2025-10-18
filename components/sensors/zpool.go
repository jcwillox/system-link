package sensors

import (
	"encoding/json"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ZpoolConfig struct {
	Pools         *[]string `yaml:"pools"`
	entity.Config `yaml:",inline"`
}

type zpoolListOutput struct {
	Pools map[string]zpoolData `json:"pools"`
}

type zpoolData struct {
	Name       string              `json:"name"`
	Properties zpoolDataProperties `json:"properties"`
}

type zpoolDataProperties struct {
	Size      zpoolProperty `json:"size"`
	Allocated zpoolProperty `json:"allocated"`
	Free      zpoolProperty `json:"free"`
}

type zpoolProperty struct {
	Value uint64 `json:"value"`
}

func getZpoolInfo(poolNames *[]string) (map[string]zpoolDataProperties, error) {
	var args []string
	if poolNames != nil {
		args = append([]string{"list", "-j", "--json-int"}, *poolNames...)
	} else {
		args = []string{"list", "-j", "--json-int"}
	}

	cmd := exec.Command("zpool", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var zpoolData zpoolListOutput
	if err := json.Unmarshal(output, &zpoolData); err != nil {
		return nil, err
	}

	result := make(map[string]zpoolDataProperties)
	for name, pool := range zpoolData.Pools {
		result[name] = pool.Properties
	}

	return result, nil
}

func NewZpool(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolsInfo, err := getZpoolInfo(cfg.Pools)
	if err != nil {
		return entities
	}

	for poolName := range poolsInfo {
		pool := poolName // capture for closure
		titleCaser := cases.Title(language.English)
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_" + pool).
			Name(titleCaser.String(pool)).
			Icon("mdi:harddisk").
			StateClass("measurement").
			Unit("%").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				poolsInfo, err := getZpoolInfo(cfg.Pools)
				if err != nil {
					return err
				}
				props, ok := poolsInfo[pool]
				if !ok {
					return nil
				}
				usedPercent := float64(props.Allocated.Value) / float64(props.Size.Value) * 100
				return e.PublishState(client, usedPercent)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewZpoolUsed(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolsInfo, err := getZpoolInfo(cfg.Pools)
	if err != nil {
		return entities
	}

	for poolName := range poolsInfo {
		pool := poolName // capture for closure
		titleCaser := cases.Title(language.English)
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_used_" + pool).
			Name(titleCaser.String(pool) + " Used").
			Icon("mdi:harddisk").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				poolsInfo, err := getZpoolInfo(cfg.Pools)
				if err != nil {
					return err
				}
				props, ok := poolsInfo[pool]
				if !ok {
					return nil
				}
				return e.PublishState(client, float64(props.Allocated.Value)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}

func NewZpoolFree(cfg ZpoolConfig) []*entity.Entity {
	var entities []*entity.Entity

	poolsInfo, err := getZpoolInfo(cfg.Pools)
	if err != nil {
		return entities
	}

	for poolName := range poolsInfo {
		pool := poolName // capture for closure
		titleCaser := cases.Title(language.English)
		newEntity := entity.NewEntity(cfg.Config).
			Type(entity.DomainSensor).
			ID("zpool_free_" + pool).
			Name(titleCaser.String(pool) + " Free").
			Icon("mdi:harddisk").
			StateClass("measurement").
			DeviceClass("data_size").
			Unit("GiB").
			Precision(1).
			Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
				poolsInfo, err := getZpoolInfo(cfg.Pools)
				if err != nil {
					return err
				}
				props, ok := poolsInfo[pool]
				if !ok {
					return nil
				}
				return e.PublishState(client, float64(props.Free.Value)/Gibibyte)
			}).Build()
		entities = append(entities, newEntity)
	}

	return entities
}
