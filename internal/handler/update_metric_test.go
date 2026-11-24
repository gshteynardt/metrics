package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type fakeService struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (f *fakeService) UpdateGauge(name string, value float64) {
	if f.gauges == nil {
		f.gauges = make(map[string]float64)
	}
	f.gauges[name] = value
}

func (f *fakeService) UpdateCounter(name string, delta int64) {
	if f.counters == nil {
		f.counters = make(map[string]int64)
	}
	f.counters[name] += delta
}

func TestUpdateMetricHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		wantStatus   int
		wantGauges   map[string]float64
		wantCounters map[string]int64
	}{
		{
			name:       "valid gauge",
			method:     http.MethodPost,
			target:     "/update/gauge/Alloc/12.5",
			wantStatus: http.StatusOK,
			wantGauges: map[string]float64{"Alloc": 12.5},
		},
		{
			name:       "valid gauge integer value",
			method:     http.MethodPost,
			target:     "/update/gauge/Alloc/100",
			wantStatus: http.StatusOK,
			wantGauges: map[string]float64{"Alloc": 100},
		},
		{
			name:       "valid gauge negative value",
			method:     http.MethodPost,
			target:     "/update/gauge/Alloc/-1.5",
			wantStatus: http.StatusOK,
			wantGauges: map[string]float64{"Alloc": -1.5},
		},
		{
			name:         "valid counter",
			method:       http.MethodPost,
			target:       "/update/counter/PollCount/5",
			wantStatus:   http.StatusOK,
			wantCounters: map[string]int64{"PollCount": 5},
		},
		{
			name:         "valid counter negative delta",
			method:       http.MethodPost,
			target:       "/update/counter/PollCount/-3",
			wantStatus:   http.StatusOK,
			wantCounters: map[string]int64{"PollCount": -3},
		},
		{
			name:       "unknown metric type",
			method:     http.MethodPost,
			target:     "/update/unknown/Alloc/1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "gauge value is not a number",
			method:     http.MethodPost,
			target:     "/update/gauge/Alloc/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "counter value is float",
			method:     http.MethodPost,
			target:     "/update/counter/PollCount/1.5",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "counter value is not a number",
			method:     http.MethodPost,
			target:     "/update/counter/PollCount/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing metric name",
			method:     http.MethodPost,
			target:     "/update/gauge/",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "missing metric value",
			method:     http.MethodPost,
			target:     "/update/gauge/Alloc/",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong method",
			method:     http.MethodGet,
			target:     "/update/gauge/Alloc/1",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeService{}

			mux := chi.NewRouter()
			mux.Post("/update/{type}/{name}/{value}", UpdateMetric(service))

			request := httptest.NewRequest(tt.method, tt.target, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Equal(t, tt.wantGauges, service.gauges)
			assert.Equal(t, tt.wantCounters, service.counters)
		})
	}
}
