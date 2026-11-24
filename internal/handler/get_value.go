package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gshteynardt/metrics/internal/model"
)

func GetValue(mr MetricReader) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := model.MetricType(chi.URLParam(r, "type"))
		name := chi.URLParam(r, "name")
		var res string

		switch metricType {
		case model.Counter:
			value, ok := mr.GetCounter(name)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}

			res = strconv.FormatInt(value, 10)

		case model.Gauge:
			value, ok := mr.GetGauge(name)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}

			res = strconv.FormatFloat(value, 'f', -1, 64)

		default:
			http.Error(rw, "unknown metric type", http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(rw, res)
	}
}
