package repository

import (
	"sync"
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

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name           string
		gaugeUpdates   []gaugeUpdate
		counterUpdates []counterUpdate
		getName        string
		wantValue      float64
		wantOK         bool
	}{
		{
			name:         "existing gauge is returned",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 123.45}},
			getName:      "Alloc",
			wantValue:    123.45,
			wantOK:       true,
		},
		{
			name:         "last written value is returned",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 1}, {"Alloc", 2}},
			getName:      "Alloc",
			wantValue:    2,
			wantOK:       true,
		},
		{
			name:         "zero value is found, not treated as missing",
			gaugeUpdates: []gaugeUpdate{{"Zero", 0}},
			getName:      "Zero",
			wantValue:    0,
			wantOK:       true,
		},
		{
			name:         "negative value is returned",
			gaugeUpdates: []gaugeUpdate{{"Temp", -1.5}},
			getName:      "Temp",
			wantValue:    -1.5,
			wantOK:       true,
		},
		{
			name:      "missing in empty storage",
			getName:   "Alloc",
			wantValue: 0,
			wantOK:    false,
		},
		{
			name:         "missing when other gauges exist",
			gaugeUpdates: []gaugeUpdate{{"HeapAlloc", 1}},
			getName:      "Alloc",
			wantValue:    0,
			wantOK:       false,
		},
		{
			name:           "counter with same name is not visible as gauge",
			counterUpdates: []counterUpdate{{"Metric", 5}},
			getName:        "Metric",
			wantValue:      0,
			wantOK:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage()

			for _, g := range tt.gaugeUpdates {
				s.UpdateGauge(g.name, g.value)
			}

			value, ok := s.GetGauge(tt.getName)

			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantOK, ok)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	tests := []struct {
		name           string
		gaugeUpdates   []gaugeUpdate
		counterUpdates []counterUpdate
		getName        string
		wantValue      int64
		wantOK         bool
	}{
		{
			name:           "existing counter is returned",
			counterUpdates: []counterUpdate{{"PollCount", 5}},
			getName:        "PollCount",
			wantValue:      5,
			wantOK:         true,
		},
		{
			name:           "accumulated value is returned",
			counterUpdates: []counterUpdate{{"PollCount", 5}, {"PollCount", 5}},
			getName:        "PollCount",
			wantValue:      10,
			wantOK:         true,
		},
		{
			name:           "zero sum is found, not treated as missing",
			counterUpdates: []counterUpdate{{"PollCount", 5}, {"PollCount", -5}},
			getName:        "PollCount",
			wantValue:      0,
			wantOK:         true,
		},
		{
			name:           "negative value is returned",
			counterUpdates: []counterUpdate{{"Neg", -3}},
			getName:        "Neg",
			wantValue:      -3,
			wantOK:         true,
		},
		{
			name:      "missing in empty storage",
			getName:   "PollCount",
			wantValue: 0,
			wantOK:    false,
		},
		{
			name:           "missing when other counters exist",
			counterUpdates: []counterUpdate{{"Requests", 1}},
			getName:        "PollCount",
			wantValue:      0,
			wantOK:         false,
		},
		{
			name:         "gauge with same name is not visible as counter",
			gaugeUpdates: []gaugeUpdate{{"Metric", 1.5}},
			getName:      "Metric",
			wantValue:    0,
			wantOK:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage()

			for _, c := range tt.counterUpdates {
				s.UpdateCounter(c.name, c.delta)
			}

			value, ok := s.GetCounter(tt.getName)

			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantOK, ok)
		})
	}
}

func TestMemStorage_All(t *testing.T) {
	tests := []struct {
		name           string
		gaugeUpdates   []gaugeUpdate
		counterUpdates []counterUpdate
		wantGauge      map[string]float64
		wantCounter    map[string]int64
	}{
		{
			name:        "empty storage returns empty non-nil maps",
			wantGauge:   map[string]float64{},
			wantCounter: map[string]int64{},
		},
		{
			name:           "gauges and counters are returned",
			gaugeUpdates:   []gaugeUpdate{{"Alloc", 1.5}, {"HeapAlloc", 3}},
			counterUpdates: []counterUpdate{{"PollCount", 2}, {"PollCount", 3}},
			wantGauge:      map[string]float64{"Alloc": 1.5, "HeapAlloc": 3},
			wantCounter:    map[string]int64{"PollCount": 5},
		},
		{
			name:         "only gauges",
			gaugeUpdates: []gaugeUpdate{{"Alloc", 1.5}},
			wantGauge:    map[string]float64{"Alloc": 1.5},
			wantCounter:  map[string]int64{},
		},
		{
			name:           "only counters",
			counterUpdates: []counterUpdate{{"PollCount", 1}},
			wantGauge:      map[string]float64{},
			wantCounter:    map[string]int64{"PollCount": 1},
		},
		{
			name:           "same name in both types is kept separately",
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

			gauge, counter := s.ListMetrics()
			assert.Equal(t, tt.wantGauge, gauge)
			assert.Equal(t, tt.wantCounter, counter)
		})
	}
}

func TestMemStorage_ConcurrentCounterUpdates(t *testing.T) {
	const (
		goroutines = 100
		iterations = 100
	)

	s := NewMemStorage()

	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			for range iterations {
				s.UpdateCounter("c", 1)
			}
		})
	}
	wg.Wait()

	value, ok := s.GetCounter("c")
	assert.True(t, ok)
	assert.Equal(t, int64(goroutines*iterations), value)
}

func TestMemStorage_ConcurrentReadWrite(t *testing.T) {
	const (
		writers    = 50
		readers    = 50
		iterations = 100
	)

	s := NewMemStorage()

	var wg sync.WaitGroup

	for range writers {
		wg.Go(func() {
			for i := range iterations {
				s.UpdateGauge("g", float64(i))
				s.UpdateCounter("c", 1)
			}
		})
	}

	for range readers {
		wg.Go(func() {
			for range iterations {
				s.GetGauge("g")
			}
		})
		wg.Go(func() {
			for range iterations {
				s.GetCounter("c")
			}
		})
		wg.Go(func() {
			for range iterations {
				gauge, counter := s.ListMetrics()
				for range gauge {
				}
				for range counter {
				}
			}
		})
	}

	wg.Wait()

	value, ok := s.GetCounter("c")
	assert.True(t, ok)
	assert.Equal(t, int64(writers*iterations), value)
}
