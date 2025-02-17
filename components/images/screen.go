package images

import (
	"bytes"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-link/entity"
	"github.com/kbinani/screenshot"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"strconv"
	"time"
)

type ScreenEntities struct {
	Timing   *entity.Config `yaml:"timing,omitempty"`
	Interval *entity.Config `yaml:"interval,omitempty"`
}

type ScreenConfig struct {
	Entities      ScreenEntities `yaml:"entities,omitempty"`
	entity.Config `yaml:",inline"`
}

func NewScreen(cfg ScreenConfig) []*entity.Entity {
	var timingEntity *entity.Entity
	var intervalEntity *entity.Entity
	var entities []*entity.Entity

	if cfg.Entities.Timing != nil {
		timingEntity = entity.NewEntity(entity.Config{}).
			Type(entity.DomainSensor).
			ID("screen_timing").
			Unit("ms").
			DeviceClass("duration").
			StateClass("measurement").
			Precision(2).
			DefaultStateTopic().
			Build()
		entities = append(entities, timingEntity)
	}

	imageEntity := entity.NewEntity(cfg.Config).
		Type(entity.DomainImage).
		ID("screen").
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
			start := time.Now()

			n := screenshot.NumActiveDisplays()
			if n == 0 {
				log.Warn().Msg("no active displays found")
				return nil
			}

			bounds := screenshot.GetDisplayBounds(0)
			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				return err
			}

			buff := new(bytes.Buffer)
			err = jpeg.Encode(buff, img, &jpeg.Options{Quality: 50})
			if err != nil {
				return err
			}

			err = e.PublishRawState(client, buff.Bytes())
			if err != nil {
				return err
			}

			if timingEntity != nil {
				err = timingEntity.PublishState(client, time.Since(start).Seconds()*1000)
				if err != nil {
					return err
				}
			}

			return nil
		}).Build()

	if cfg.Entities.Interval != nil {
		intervalEntity = entity.NewEntity(*cfg.Entities.Interval).
			Type(entity.DomainNumber).
			ID("screen_interval").
			Unit("s").
			Min(0.1).
			Max(300).
			Step(0.1).
			EntityCategory("config").
			OnCommand(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
				_ = e.PublishState(client, string(message.Payload()))
			}).
			OnState(func(entity *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler, message mqtt.Message) {
				interval, err := strconv.ParseFloat(string(message.Payload()), 64)
				if err != nil {
					log.Err(err).Msg("failed to parse interval")
					return
				}
				interval *= float64(time.Second)
				err = imageEntity.UpdateJob(scheduler, gocron.DurationJob(time.Duration(interval)))
				if err != nil {
					log.Err(err).Msg("failed to set interval")
				}
			}).Build()
		entities = append(entities, intervalEntity)
	}

	return append(entities, imageEntity)
}
