package engine

import (
	"time"
)

type Filters struct {
	Or  []Filters `yaml:"or"`
	And []Filters `yaml:"and"`

	Throttle time.Duration `yaml:"throttle"`
	Delta    float64       `yaml:"delta"`

	state *FilterState
}

type FilterState struct {
	LastState   interface{}
	LastUpdated time.Time
}

// Filter returns true if the state passes the filters
func (f Filters) Filter(state interface{}) (interface{}, error) {
	state, err := f.filter(state)
	if err != nil {
		f.state.LastState = state
		f.state.LastUpdated = time.Now()
		return state, err
	} else {
		return nil, nil
	}
}

func (f Filters) filter(state interface{}) (interface{}, error) {
	var err error

	for _, filters := range f.Or {
		filters.state = f.state

		// if any filter passes continue, otherwise fail
		nextState, err := filters.filter(state)
		if err != nil {
			return nextState, err
		}
		if nextState != nil {
			state = nextState
			break
		}
	}

	if state, err = f.throttle(state); state == nil {
		return nil, err
	}

	if state, err = f.delta(state); state == nil {
		return nil, err
	}

	return state, nil
}

func (f Filters) throttle(state interface{}) (interface{}, error) {
	if f.Throttle > 0 {
		// check if throttle has passed
		if time.Since(f.state.LastUpdated) < f.Throttle {
			return nil, nil
		}
	}
	return state, nil
}

func (f Filters) delta(state interface{}) (interface{}, error) {
	if f.Delta != 0 {
		// check if delta has passed
		if f.state.LastState != nil {
			switch f.state.LastState.(type) {
			case float64:
				if f.Delta > 0 {
					if state.(float64) < f.state.LastState.(float64)-f.Delta {
						return nil, nil
					}
				} else {
					if state.(float64) > f.state.LastState.(float64)-f.Delta {
						return nil, nil
					}
				}
			case int:
				if f.Delta > 0 {
					if state.(int) < f.state.LastState.(int)-int(f.Delta) {
						return nil, nil
					}
				} else {
					if state.(int) > f.state.LastState.(int)-int(f.Delta) {
						return nil, nil
					}
				}
			}
		}
	}
	return state, nil
}
