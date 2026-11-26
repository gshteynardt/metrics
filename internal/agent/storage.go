package agent

import (
	"maps"
	"sync"
)

// Storage is a concurrency-safe in-memory buffer for metrics collected by the agent.
type Storage struct {
	mu        sync.Mutex
	gauges    map[string]float64
	pollCount int64
}

// NewStorage returns an empty Storage.
func NewStorage() *Storage {
	return &Storage{
		gauges: make(map[string]float64),
	}
}

// Update merges gauges and increments PollCount.
func (s *Storage) Update(gauges map[string]float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	maps.Copy(s.gauges, gauges)
	s.pollCount++
}

// Snapshot returns a copy of the stored gauges and the current PollCount.
func (s *Storage) Snapshot() (map[string]float64, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return maps.Clone(s.gauges), s.pollCount
}
