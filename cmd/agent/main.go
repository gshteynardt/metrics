package main

import (
	"log"
	"time"

	"github.com/gshteynardt/metrics/internal/agent"
)

func main() {
	run()
}

const (
	serverAddr     = "http://localhost:8080"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func run() {
	storage := agent.NewStorage()
	collector := agent.NewCollector(storage)
	sender := agent.NewSender(storage, serverAddr)

	go func() {
		for {
			time.Sleep(reportInterval)
			if err := sender.Send(); err != nil {
				log.Printf("report failed: %v", err)
			} else {
				log.Println("metrics sent")
			}
		}
	}()

	for {
		collector.Collect()
		log.Println("metrics collected")

		time.Sleep(pollInterval)
	}
}
