package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gshteynardt/metrics/internal/model"
)

// UpdateMetric returns a handler for POST /update/{type}/{name}/{value}.
// Gauges are replaced and counters are incremented by the given delta.
func UpdateMetric(w MetricWriter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := model.MetricType(chi.URLParam(r, "type"))
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		switch metricType {
		case model.Counter:
			delta, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(rw, "invalid counter value", http.StatusBadRequest)
				return
			}

			w.UpdateCounter(name, delta)

		case model.Gauge:
			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(rw, "invalid gauge value", http.StatusBadRequest)
				return
			}

			w.UpdateGauge(name, number)

		default:
			http.Error(rw, "unknown metric type", http.StatusBadRequest)
			return
		}
	}
}
