package repository

import "sync"

type MemStorage struct {
	gauge           map[string]float64
	counter         map[string]int64
	gaugeMu   sync.Mutex
	counterMu sync.Mutex
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
