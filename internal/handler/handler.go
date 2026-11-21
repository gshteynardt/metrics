// Package handler implements HTTP handlers for metric updates.
package handler

import "net/http"

type MetricService interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, delta int64)
}

type Handler struct {
	service MetricService
}

func NewHandler(service MetricService) *Handler {
	return &Handler{service: service}
}

// Routes returns an http.Handler with all metric endpoints registered.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", h.UpdateMetric)

	return mux
}
