package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMetricsHandler(t *testing.T) {
	tests := []struct {
		name        string
		reader      *fakeReader
		wantInOrder []string
		wantAbsent  []string
	}{
		{
			name:   "empty storage renders page",
			reader: &fakeReader{},
		},
		{
			name: "metrics are rendered sorted, gauges before counters",
			reader: &fakeReader{
				gauges: map[string]float64{
					"HeapAlloc": 3, "Alloc": 1.5, "GCSys": 4, "BuckHashSys": 2, "Frees": 5,
				},
				counters: map[string]int64{"PollCount": 10, "Errors": 1},
			},
			wantInOrder: []string{
				"<td>Alloc</td><td>gauge</td><td>1.5</td>",
				"<td>BuckHashSys</td><td>gauge</td><td>2</td>",
				"<td>Frees</td><td>gauge</td><td>5</td>",
				"<td>GCSys</td><td>gauge</td><td>4</td>",
				"<td>HeapAlloc</td><td>gauge</td><td>3</td>",
				"<td>Errors</td><td>counter</td><td>1</td>",
				"<td>PollCount</td><td>counter</td><td>10</td>",
			},
		},
		{
			name: "metric name is html-escaped",
			reader: &fakeReader{
				gauges: map[string]float64{"<script>": 1},
			},
			wantInOrder: []string{"<td>&lt;script&gt;</td>"},
			wantAbsent:  []string{"<script>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			ListMetrics(tt.reader).ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))

			raw, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			body := string(raw)

			rest := body
			for _, s := range tt.wantInOrder {
				i := strings.Index(rest, s)
				require.GreaterOrEqual(t, i, 0, "%q not found (or out of order) in body:\n%s", s, body)
				rest = rest[i+len(s):]
			}

			for _, s := range tt.wantAbsent {
				assert.NotContains(t, body, s)
			}
		})
	}
}
