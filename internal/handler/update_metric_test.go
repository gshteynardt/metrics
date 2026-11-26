package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			store := &fakeStore{}
			rec := serve(t, store, tt.method, tt.target)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantGauges, store.gauges)
			assert.Equal(t, tt.wantCounters, store.counters)
		})
	}
}
