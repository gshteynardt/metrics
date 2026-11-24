package handler

type fakeReader struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (f *fakeReader) GetGauge(name string) (float64, bool) {
	v, ok := f.gauges[name]
	return v, ok
}

func (f *fakeReader) GetCounter(name string) (int64, bool) {
	v, ok := f.counters[name]
	return v, ok
}

func (f *fakeReader) ListMetrics() (map[string]float64, map[string]int64) {
	return f.gauges, f.counters
}
