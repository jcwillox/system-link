package sensors

import (
	"encoding/json"
	"github.com/PaesslerAG/jsonpath"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/jcwillox/system-link/utils"
	"github.com/rs/zerolog/log"
)

type CustomConfig struct {
	utils.CommandConfig `yaml:",inline"`
	entity.Config       `yaml:",inline"`
	JsonAttributes      []string `yaml:"json_attributes,omitempty"`
	JsonAttributesPath  string   `yaml:"json_attributes_path,omitempty"`
}

func NewCustom(cfg CustomConfig) *entity.Entity {
	return entity.NewEntity(cfg.Config).
		Type(entity.DomainSensor).
		ObjectID(cfg.UniqueID).
		DefaultAttributesTopic().
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			res, err := utils.RunCommand(cfg.CommandConfig)
			if err != nil {
				log.Err(err).Str("command", cfg.CommandConfig.Command).Msg("failed to run command")
			}

			// publish stdout as state
			err = e.PublishRawState(client, res.Stdout)
			if err != nil {
				return err
			}

			// check if we need to publish attributes
			if len(cfg.JsonAttributes) == 0 {
				return nil
			}

			// try parse as json for attributes
			var attributesBase map[string]interface{}
			err = json.Unmarshal(res.Stdout, &attributesBase)
			if err != nil {
				log.Err(err).Msg("failed to parse stdout as json")
			}

			// handle attributes path
			if cfg.JsonAttributesPath != "" {
				attributes, err := jsonpath.Get(cfg.JsonAttributesPath, attributesBase)
				if err != nil {
					log.Err(err).Str("path", cfg.JsonAttributesPath).Msg("failed to get json path")
				}
				if v, ok := attributes.(map[string]interface{}); ok {
					attributesBase = v
				} else {
					log.Error().Str("path", cfg.JsonAttributesPath).Msg("json path did not return object")
				}
			}

			// gather attributes
			attributes := make(map[string]interface{})
			if len(cfg.JsonAttributes) > 0 {
				for _, attribute := range cfg.JsonAttributes {
					attributes[attribute] = attributesBase[attribute]
				}
			}

			// publish attributes
			err = e.PublishAttributes(client, attributes)
			if err != nil {
				return err
			}

			return nil
		}).Build()
}
