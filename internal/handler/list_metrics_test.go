package handler

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMetricsHandler(t *testing.T) {
	tests := []struct {
		name        string
		store       *fakeStore
		wantInOrder []string
		wantAbsent  []string
	}{
		{
			name:  "empty storage renders page",
			store: &fakeStore{},
		},
		{
			name: "metrics are rendered sorted, gauges before counters",
			store: &fakeStore{
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
			store: &fakeStore{
				gauges: map[string]float64{"<script>": 1},
			},
			wantInOrder: []string{"<td>&lt;script&gt;</td>"},
			wantAbsent:  []string{"<script>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, tt.store, http.MethodGet, "/")

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

			body := rec.Body.String()

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
