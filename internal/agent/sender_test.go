package agent

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeStorage stores metric and returns snapshot of metric
type fakeStorage struct {
	gauges    map[string]float64
	pollCount int64
}

func (f *fakeStorage) Snapshot() (map[string]float64, int64) {
	return maps.Clone(f.gauges), f.pollCount
}

func TestSender_Send(t *testing.T) {
	tests := []struct {
		name         string
		gauges       map[string]float64
		pollCount    int64
		serverStatus int
		wantRequests []string
		wantErr      bool
	}{
		{
			name:         "single gauge and poll count",
			gauges:       map[string]float64{"Alloc": 1.5},
			pollCount:    3,
			serverStatus: http.StatusOK,
			wantRequests: []string{
				"POST /update/gauge/Alloc/1.5",
				"POST /update/counter/PollCount/3",
			},
			wantErr: false,
		},
		{
			name:         "multiple gauges",
			gauges:       map[string]float64{"Alloc": 1.5, "HeapAlloc": 3, "RandomValue": 0.42},
			pollCount:    5,
			serverStatus: http.StatusOK,
			wantRequests: []string{
				"POST /update/gauge/Alloc/1.5",
				"POST /update/gauge/HeapAlloc/3",
				"POST /update/gauge/RandomValue/0.42",
				"POST /update/counter/PollCount/5",
			},
			wantErr: false,
		},
		{
			name:         "gauge values are formatted without exponent",
			gauges:       map[string]float64{"Big": 1000000, "Small": 0.000001},
			pollCount:    1,
			serverStatus: http.StatusOK,
			wantRequests: []string{
				"POST /update/gauge/Big/1000000",
				"POST /update/gauge/Small/0.000001",
				"POST /update/counter/PollCount/1",
			},
			wantErr: false,
		},
		{
			name:         "empty snapshot sends only poll count",
			gauges:       map[string]float64{},
			pollCount:    0,
			serverStatus: http.StatusOK,
			wantRequests: []string{
				"POST /update/counter/PollCount/0",
			},
			wantErr: false,
		},
		{
			name:         "server error stops sending",
			gauges:       map[string]float64{"Alloc": 1.5},
			pollCount:    3,
			serverStatus: http.StatusInternalServerError,
			wantRequests: []string{
				"POST /update/gauge/Alloc/1.5",
			},
			wantErr: true,
		},
		{
			name:         "bad request is an error",
			gauges:       map[string]float64{"Alloc": 1.5},
			pollCount:    3,
			serverStatus: http.StatusBadRequest,
			wantRequests: []string{
				"POST /update/gauge/Alloc/1.5",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = append(got, r.Method+" "+r.URL.Path)
				assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
				w.WriteHeader(tt.serverStatus)
			}))

			defer srv.Close()
			fake := &fakeStorage{gauges: tt.gauges, pollCount: tt.pollCount}
			err := NewSender(fake, srv.URL).Send()
			assert.ElementsMatch(t, tt.wantRequests, got)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSender_Send_ServerUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	addr := srv.URL
	srv.Close()

	fake := &fakeStorage{gauges: map[string]float64{"Alloc": 1.5}, pollCount: 1}

	err := NewSender(fake, addr).Send()

	assert.Error(t, err)
}
