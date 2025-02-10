package filters

import (
	"fmt"
	"math"
	"time"
)

type Filter struct {
	*OrFilter       `yaml:",inline"`
	*ThrottleFilter `yaml:",inline"`
	*DeltaFilter    `yaml:",inline"`
}

type FilterI interface {
	Filter(state interface{}) (interface{}, error)
	IsSet() bool
	fmt.Stringer
}

type OrFilter struct {
	Or []Filter `yaml:"or"`
}

func (f *OrFilter) String() string {
	return fmt.Sprintf("OrFilter{Or: %v}", f.Or)
}

func (f *OrFilter) IsSet() bool {
	return len(f.Or) > 0
}

func (f *OrFilter) Filter(state interface{}) (interface{}, error) {
	filters := make([]FilterI, 0, len(f.Or)*3)
	for _, filter := range f.Or {
		filters = append(filters, filter.OrFilter, filter.ThrottleFilter, filter.DeltaFilter)
	}
	// run through all filters
	var err error
	for _, filter := range filters {
		if !filter.IsSet() {
			continue
		}
		// find first filter that passes
		var newState interface{}
		newState, err = filter.Filter(state)
		if err == nil {
			return newState, nil
		}
	}
	return state, err
}

type ThrottleFilter struct {
	Throttle time.Duration `yaml:"throttle"`
	since    time.Time
}

func (f *ThrottleFilter) String() string {
	return fmt.Sprintf("ThrottleFilter{Throttle: %v, since: %v}", f.Throttle, f.since)
}

func (f *ThrottleFilter) IsSet() bool {
	return f.Throttle > 0
}

func (f *ThrottleFilter) Filter(state interface{}) (interface{}, error) {
	if f.since.IsZero() {
		f.since = time.Now()
	}
	if time.Since(f.since) < f.Throttle {
		return nil, SkipSendErr
	}
	// time we last allowed a state to pass
	f.since = time.Now()
	return state, nil
}

type DeltaFilter struct {
	Delta float64 `yaml:"delta"`
	state interface{}
}

func (f *DeltaFilter) String() string {
	return fmt.Sprintf("DeltaFilter{Delta: %v, state: %v}", f.Delta, f.state)
}

func (f *DeltaFilter) IsSet() bool {
	return f.Delta > 0
}

func (f *DeltaFilter) Filter(state interface{}) (interface{}, error) {
	// check if delta has passed
	if f.state != nil {
		switch value := f.state.(type) {
		case float64:
			diff := math.Abs(value - state.(float64))
			if diff < f.Delta {
				return state, SkipSendErr
			}
		case int:
			diff := math.Abs(float64(value - state.(int)))
			if diff < f.Delta {
				return state, SkipSendErr
			}
		}
	}
	f.state = state
	return state, nil
}
