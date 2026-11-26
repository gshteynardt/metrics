package handler

import (
	"bytes"
	"html/template"
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/gshteynardt/metrics/internal/model"
)

type metricRow struct {
	Name  string
	Type  model.MetricType
	Value string
}

var listTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Metrics</title></head>
<body>
<h1>Metrics</h1>
<table>
<tr><th>Name</th><th>Type</th><th>Value</th></tr>
{{range .}}<tr><td>{{.Name}}</td><td>{{.Type}}</td><td>{{.Value}}</td></tr>
{{end}}</table>
</body>
</html>`))

// ListMetrics returns a handler for GET / that renders all metrics as an HTML table.
func ListMetrics(mr MetricReader) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		gauges, counters := mr.ListMetrics()

		rows := make([]metricRow, 0, len(gauges)+len(counters))

		for _, name := range slices.Sorted(maps.Keys(gauges)) {
			rows = append(rows, metricRow{
				Name:  name,
				Type:  model.Gauge,
				Value: strconv.FormatFloat(gauges[name], 'f', -1, 64),
			})
		}

		for _, name := range slices.Sorted(maps.Keys(counters)) {
			rows = append(rows, metricRow{
				Name:  name,
				Type:  model.Counter,
				Value: strconv.FormatInt(counters[name], 10),
			})
		}

		var buf bytes.Buffer
		if err := listTmpl.Execute(&buf, rows); err != nil {
			http.Error(rw, "failed to render page", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(rw)
	}
}
