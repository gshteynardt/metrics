package handler

import (
	"net/http"
	"strconv"

	"github.com/gshteynardt/metrics/internal/model"
)

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := model.MetricType(r.PathValue("type"))
	name := r.PathValue("name")
	value := r.PathValue("value")

	switch metricType {
	case model.Counter:
		delta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}

		h.service.UpdateCounter(name, delta)

	case model.Gauge:
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}

		h.service.UpdateGauge(name, number)

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
