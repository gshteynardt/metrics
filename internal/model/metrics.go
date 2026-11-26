// Package model defines metric types shared by the agent and the server.
package model

// MetricType is the kind of a metric: gauge or counter.
type MetricType string

const (
	// Counter is a metric whose value accumulates deltas.
	Counter MetricType = "counter"
	// Gauge is a metric whose value is replaced on each update.
	Gauge MetricType = "gauge"
)

// Metrics is the JSON representation of a single metric.
// Delta is set for counters and Value for gauges. Both are pointers
// to distinguish a zero value from an absent one, so that an unset
// field is omitted from the encoded JSON.
type Metrics struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
	Hash  string     `json:"hash,omitempty"`
}
