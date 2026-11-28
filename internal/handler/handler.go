// Package handler implements HTTP handlers for reading and updating metrics.
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MetricWriter updates gauge and counter metrics.
type MetricWriter interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, delta int64)
}

// MetricReader provides read access to stored metrics.
type MetricReader interface {
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ListMetrics() (map[string]float64, map[string]int64)
}

// MetricStorage is the metric storage used by the handlers.
type MetricStorage interface {
	MetricWriter
	MetricReader
}

// Routes returns an http.Handler with all metric endpoints registered.
func Routes(s MetricStorage) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", UpdateMetric(s))
	r.Get("/value/{type}/{name}", GetValue(s))
	r.Get("/", ListMetrics(s))

	return r
}
