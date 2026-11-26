package handler

type fakeStore struct {
	gauges   map[string]float64
	counters map[string]int64
}

var _ MetricStorage = (*fakeStore)(nil)

func (f *fakeStore) GetGauge(name string) (float64, bool) {
	v, ok := f.gauges[name]
	return v, ok
}

func (f *fakeStore) GetCounter(name string) (int64, bool) {
	v, ok := f.counters[name]
	return v, ok
}

func (f *fakeStore) ListMetrics() (map[string]float64, map[string]int64) {
	return f.gauges, f.counters
}

func (f *fakeStore) UpdateGauge(name string, value float64) {
	if f.gauges == nil {
		f.gauges = make(map[string]float64)
	}
	f.gauges[name] = value
}

func (f *fakeStore) UpdateCounter(name string, delta int64) {
	if f.counters == nil {
		f.counters = make(map[string]int64)
	}
	f.counters[name] += delta
}
