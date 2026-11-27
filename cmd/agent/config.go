package main

import (
	"flag"
	"fmt"
	"time"
)

type config struct {
	serverURL      string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func parseFlags(args []string) (config, error) {
	var (
		addr   string
		report int
		poll   int
	)

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.StringVar(&addr, "a", "localhost:8080", "HTTP server address (host:port)")
	fs.IntVar(&report, "r", 10, "interval between sending metrics to the server, in `seconds`")
	fs.IntVar(&poll, "p", 2, "interval between collecting runtime metrics, in `seconds`")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: agent [flags]\n\n")
		fmt.Fprintf(fs.Output(), "Collects runtime metrics and sends them to the metrics server.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	if fs.NArg() > 0 {
		return config{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	if report <= 0 {
		return config{}, fmt.Errorf("report interval must be greater than 0, got %d", report)
	}

	if poll <= 0 {
		return config{}, fmt.Errorf("poll interval must be greater than 0, got %d", poll)
	}

	return config{
		serverURL:      "http://" + addr,
		reportInterval: time.Duration(report) * time.Second,
		pollInterval:   time.Duration(poll) * time.Second,
	}, nil
}
