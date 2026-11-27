package main

import (
	"log"
	"os"
	"time"

	"github.com/gshteynardt/metrics/internal/agent"
)

func main() {
	run()
}

func run() {
	cfg, err := parseFlags(os.Args[1:])

	if err != nil {
		log.Fatal(err)
	}

	storage := agent.NewStorage()
	collector := agent.NewCollector(storage)
	sender := agent.NewSender(storage, cfg.serverURL)

	go func() {
		for {
			time.Sleep(cfg.reportInterval)
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

		time.Sleep(cfg.pollInterval)
	}
}
