// Package handler implements HTTP handlers for metric updates.
package handler


type MetricService interface {
    UpdateGauge(name string, value float64)
    UpdateCounter(name string, delta int64)
}

type Handler struct {
    service MetricService
}

func New(service MetricService) *Handler {
    return &Handler{service: service}
}
