package agent

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gshteynardt/metrics/internal/model"
)

type metricsReader interface {
	Snapshot() (map[string]float64, int64)
}

type Sender struct {
	storage    metricsReader
	serverAddr string
	client     *http.Client
}

func NewSender(storage metricsReader, serverAddr string) *Sender {
	return &Sender{
		storage:    storage,
		serverAddr: serverAddr,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *Sender) sendOne(metricType, name, value string) error {
	path := fmt.Sprintf("%s/update/%s/%s/%s", s.serverAddr, metricType, name, value)
	response, err := s.client.Post(path, "text/plain", nil)

	if err != nil {
		return fmt.Errorf("send %s %s: %w", metricType, name, err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("send %s %s: unexpected status %d", metricType, name, response.StatusCode)
	}

	return nil
}

func (s *Sender) Send() error {
	gauges, pollCount := s.storage.Snapshot()

	for name, value := range gauges {
		if err := s.sendOne(string(model.Gauge), name, strconv.FormatFloat(value, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if err := s.sendOne(string(model.Counter), "PollCount", strconv.FormatInt(pollCount, 10)); err != nil {
		return err
	}

	return nil
}
