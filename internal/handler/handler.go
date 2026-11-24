// Package handler implements HTTP handlers for metric updates.
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricWriter interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, delta int64)
}

type MetricReader interface {
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ListMetrics() (map[string]float64, map[string]int64)
}

type MetricService interface {
	MetricWriter
	MetricReader
}

// implements HTTP handlers for metrics
func Routes(s MetricService) http.Handler {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", UpdateMetric(s))
	r.Get("/value/{type}/{name}", GetValue(s))
	r.Get("/", ListMetrics(s))

	return r
}
