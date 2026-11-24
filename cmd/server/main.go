package main

import (
	"net/http"

	"github.com/gshteynardt/metrics/internal/handler"
	"github.com/gshteynardt/metrics/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()

	err := http.ListenAndServe(":8080", handler.Routes(storage))
	if err != nil {
		panic(err)
	}
}
