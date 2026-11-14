package main

import (
	"net/http"

	"github.com/gshteynardt/metrics/internal/handler"
	"github.com/gshteynardt/metrics/internal/repository"
)

func main() {
	mux := http.NewServeMux()
	storage := repository.NewMemStorage()
	h := handler.New(storage)

	mux.HandleFunc("POST /update/{type}/{name}/{value}", h.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
