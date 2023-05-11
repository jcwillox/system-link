package utils

import (
	"gopkg.in/yaml.v3"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	var s string
	err := n.Decode(&s)
	if err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}
