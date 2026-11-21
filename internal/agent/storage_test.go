package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name          string
		updates       []map[string]float64
		wantGauges    map[string]float64
		wantPollCount int64
	}{
		{
			name:          "empty storage",
			updates:       nil,
			wantGauges:    map[string]float64{},
			wantPollCount: 0,
		},
		{
			name: "single update",
			updates: []map[string]float64{
				{"Alloc": 1.5, "HeapAlloc": 3},
			},
			wantGauges:    map[string]float64{"Alloc": 1.5, "HeapAlloc": 3},
			wantPollCount: 1,
		},
		{
			name: "existing gauge is replaced",
			updates: []map[string]float64{
				{"Alloc": 1.5},
				{"Alloc": 2.5},
			},
			wantGauges:    map[string]float64{"Alloc": 2.5},
			wantPollCount: 2,
		},
		{
			name: "gauges from different updates are merged",
			updates: []map[string]float64{
				{"Alloc": 1.5},
				{"RandomValue": 0.42},
			},
			wantGauges:    map[string]float64{"Alloc": 1.5, "RandomValue": 0.42},
			wantPollCount: 2,
		},
		{
			name: "empty update still increments poll count",
			updates: []map[string]float64{
				{"Alloc": 1.5},
				{},
			},
			wantGauges:    map[string]float64{"Alloc": 1.5},
			wantPollCount: 2,
		},
		{
			name: "poll count grows with every update",
			updates: []map[string]float64{
				{"Alloc": 1},
				{"Alloc": 2},
				{"Alloc": 3},
				{"Alloc": 4},
				{"Alloc": 5},
			},
			wantGauges:    map[string]float64{"Alloc": 5},
			wantPollCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStorage()

			for _, u := range tt.updates {
				s.Update(u)
			}

			gauges, pollCount := s.Snapshot()
			assert.Equal(t, tt.wantGauges, gauges)
			assert.Equal(t, tt.wantPollCount, pollCount)
		})
	}
}

func TestStorage_SnapshotReturnsCopy(t *testing.T) {
	s := NewStorage()
	s.Update(map[string]float64{"Alloc": 1})

	snapshot, _ := s.Snapshot()
	snapshot["Alloc"] = 100
	snapshot["Injected"] = 42

	got, _ := s.Snapshot()
	assert.Equal(t, map[string]float64{"Alloc": 1}, got)
}
