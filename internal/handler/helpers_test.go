package handler

import (
	"net/http/httptest"
	"testing"
)

func serve(t *testing.T, store MetricStorage, method, target string) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	Routes(store).ServeHTTP(rec, httptest.NewRequest(method, target, nil))

	return rec
}
