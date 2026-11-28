package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValueHandler(t *testing.T) {
	store := &fakeStore{
		gauges:   map[string]float64{"Alloc": 123.45},
		counters: map[string]int64{"PollCount": 10, "Shared": 1},
	}

	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "gauge found",
			target:     "/value/gauge/Alloc",
			wantStatus: http.StatusOK,
			wantBody:   "123.45",
		},
		{
			name:       "counter found",
			target:     "/value/counter/PollCount",
			wantStatus: http.StatusOK,
			wantBody:   "10",
		},
		{
			name:       "gauge not found (exists only as counter)",
			target:     "/value/gauge/Shared",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "counter not found (exists only as gauge)",
			target:     "/value/counter/Alloc",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "unknown metric type",
			target:     "/value/unknown/Alloc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, store, http.MethodGet, tt.target)
			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus != http.StatusOK {
				return
			}

			assert.Equal(t, tt.wantBody, rec.Body.String())
			assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
		})
	}
}
