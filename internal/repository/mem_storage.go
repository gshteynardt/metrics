package repository

import (
	"maps"
	"sync"
)

type MemStorage struct {
	gauge     map[string]float64
	counter   map[string]int64
	gaugeMu   sync.RWMutex
	counterMu sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.gaugeMu.Lock()
	s.gauge[name] = value
	s.gaugeMu.Unlock()
}

func (s *MemStorage) UpdateCounter(name string, delta int64) {
	s.counterMu.Lock()
	s.counter[name] += delta
	s.counterMu.Unlock()
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	s.gaugeMu.RLock()
	value, ok := s.gauge[name]
	s.gaugeMu.RUnlock()
	return value, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	s.counterMu.RLock()
	value, ok := s.counter[name]
	s.counterMu.RUnlock()
	return value, ok
}

func (s *MemStorage) ListMetrics() (map[string]float64, map[string]int64) {
	s.gaugeMu.RLock()
	gauges := maps.Clone(s.gauge)
	s.gaugeMu.RUnlock()

	s.counterMu.RLock()
	counters := maps.Clone(s.counter)
	s.counterMu.RUnlock()

	return gauges, counters
}
