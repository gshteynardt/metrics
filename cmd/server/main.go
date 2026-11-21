package main

import (
	"net/http"

	"github.com/gshteynardt/metrics/internal/handler"
	"github.com/gshteynardt/metrics/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)

	err := http.ListenAndServe(":8080", h.Routes())
	if err != nil {
		panic(err)
	}
}
