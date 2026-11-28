package agent

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// expectedGauges lists every gauge the agent must report according to the task.
// It is written out explicitly instead of being derived from runtimeGauges,
// so the test does not verify the collector against itself.
var expectedGauges = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
}

func TestCollector_Collect(t *testing.T) {
	tests := []struct {
		name          string
		collectCalls  int
		wantGauges    []string
		wantPollCount int64
	}{
		{
			name:          "no collect calls",
			collectCalls:  0,
			wantGauges:    nil,
			wantPollCount: 0,
		},
		{
			name:          "single collect",
			collectCalls:  1,
			wantGauges:    expectedGauges,
			wantPollCount: 1,
		},
		{
			name:          "multiple collects",
			collectCalls:  3,
			wantGauges:    expectedGauges,
			wantPollCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStorage()
			c := NewCollector(s)

			for range tt.collectCalls {
				c.Collect()
			}

			gauges, pollCount := s.Snapshot()
			keys := slices.Collect(maps.Keys(gauges))
			assert.Equal(t, tt.wantPollCount, pollCount)
			assert.ElementsMatch(t, tt.wantGauges, keys)
		})
	}
}
