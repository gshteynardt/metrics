// Package repository provides metric storage implementations.
package repository

import (
	"maps"
	"sync"
)

// MemStorage is a concurrency-safe in-memory metric storage.
type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage returns an empty MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateGauge sets the gauge to value, replacing the previous one.
func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gauges[name] = value
}

// UpdateCounter adds delta to the counter.
func (s *MemStorage) UpdateCounter(name string, delta int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counters[name] += delta
}

// GetGauge returns the gauge value and whether it exists.
func (s *MemStorage) GetGauge(name string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.gauges[name]

	return value, ok
}

// GetCounter returns the counter value and whether it exists.
func (s *MemStorage) GetCounter(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.counters[name]

	return value, ok
}

// ListMetrics returns copies of all gauges and counters taken at the same moment.
func (s *MemStorage) ListMetrics() (map[string]float64, map[string]int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return maps.Clone(s.gauges), maps.Clone(s.counters)
}
