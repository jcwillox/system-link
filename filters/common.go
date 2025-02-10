package filters

import (
	"errors"
)

var SkipSendErr = errors.New("skip sending state")

type Filters struct {
	Filters []Filter `yaml:"filters"`
}

func (f *Filters) Filter(state interface{}) (interface{}, error) {
	filters := make([]FilterI, 0, len(f.Filters)*3)
	for _, filter := range f.Filters {
		filters = append(filters, filter.OrFilter, filter.ThrottleFilter, filter.DeltaFilter)
	}
	// run through all filters
	var err error
	for _, filter := range filters {
		if !filter.IsSet() {
			continue
		}
		// pass output state to next filter
		state, err = filter.Filter(state)
		// ensure all filters pass
		if err != nil {
			return state, err
		}
	}
	return state, nil
}
