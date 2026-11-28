package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gshteynardt/metrics/internal/handler"
	"github.com/gshteynardt/metrics/internal/repository"
)

func main() {
	run()
}

func run() {
	cfg, parseErr := parseFlags(os.Args[1:])

	if parseErr != nil {
		log.Fatal(parseErr)
	}

	storage := repository.NewMemStorage()

	err := http.ListenAndServe(cfg.addr, handler.Routes(storage))
	if err != nil {
		panic(err)
	}
}
