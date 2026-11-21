package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type gaugeUpdate struct {
	name  string
	value float64
}

type counterUpdate struct {
	name  string
	delta int64
}

func TestMemStorage(t *testing.T) {
	tests := []struct {
		name           string
		gaugeUpdates   []gaugeUpdate
		counterUpdates []counterUpdate
		wantGauge      map[string]float64
		wantCounter    map[string]int64
	}{
		{
			name:         "gauge is set",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 1.5}},
			wantGauge:    map[string]float64{"Alloc": 1.5},
			wantCounter:  map[string]int64{},
		},
		{
			name:         "gauge is replaced",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 1.5}, {"Alloc", 2.5}},
			wantGauge:    map[string]float64{"Alloc": 2.5},
			wantCounter:  map[string]int64{},
		},
		{
			name:         "different gauges are kept separately",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 1.5}, {"HeapAlloc", 3}},
			wantGauge:    map[string]float64{"Alloc": 1.5, "HeapAlloc": 3},
			wantCounter:  map[string]int64{},
		},
		{
			name:           "counter is set",
			counterUpdates: []counterUpdate{{"PollCount", 3}},
			wantGauge:      map[string]float64{},
			wantCounter:    map[string]int64{"PollCount": 3},
		},
		{
			name:           "counter is accumulated",
			counterUpdates: []counterUpdate{{"PollCount", 3}, {"PollCount", 4}},
			wantGauge:      map[string]float64{},
			wantCounter:    map[string]int64{"PollCount": 7},
		},
		{
			name:           "counter accepts negative delta",
			counterUpdates: []counterUpdate{{"PollCount", 5}, {"PollCount", -2}},
			wantGauge:      map[string]float64{},
			wantCounter:    map[string]int64{"PollCount": 3},
		},
		{
			name:           "different counters are kept separately",
			counterUpdates: []counterUpdate{{"PollCount", 1}, {"Requests", 10}},
			wantGauge:      map[string]float64{},
			wantCounter:    map[string]int64{"PollCount": 1, "Requests": 10},
		},
		{
			name:           "gauge and counter with same name do not mix",
			gaugeUpdates:   []gaugeUpdate{{"Metric", 1.5}},
			counterUpdates: []counterUpdate{{"Metric", 2}},
			wantGauge:      map[string]float64{"Metric": 1.5},
			wantCounter:    map[string]int64{"Metric": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage()

			for _, g := range tt.gaugeUpdates {
				s.UpdateGauge(g.name, g.value)
			}

			for _, c := range tt.counterUpdates {
				s.UpdateCounter(c.name, c.delta)
			}

			assert.Equal(t, tt.wantGauge, s.gauge)
			assert.Equal(t, tt.wantCounter, s.counter)
		})
	}
}
