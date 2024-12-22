package images

import (
	"bytes"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron/v2"
	"github.com/jcwillox/system-bridge/entity"
	"github.com/kbinani/screenshot"
	"github.com/rs/zerolog/log"
	"image/jpeg"
)

func NewScreen(cfg entity.Config) *entity.Entity {
	return entity.NewEntity(cfg).
		Type(entity.DomainImage).
		ID("screen").
		Schedule(func(e *entity.Entity, client mqtt.Client, scheduler gocron.Scheduler) error {
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

			return e.PublishRawState(client, buff.Bytes())
		}).Build()
}
