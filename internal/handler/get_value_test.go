package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValueHandler(t *testing.T) {
	reader := &fakeReader{
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
			mux := chi.NewRouter()
			mux.Get("/value/{type}/{name}", GetValue(reader))

			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantStatus != http.StatusOK {
				return
			}

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantBody, string(body))
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
		})
	}
}
