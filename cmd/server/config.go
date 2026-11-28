package main

import (
	"flag"
	"fmt"
)

type config struct {
	addr string
}

func parseFlags(args []string) (config, error) {
	var addr string

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&addr, "a", "localhost:8080", "`address` to listen on (host:port)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: server [flags]\n\n")
		fmt.Fprintf(fs.Output(), "Accepts metrics over HTTP and stores them in memory.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	if fs.NArg() > 0 {
		return config{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	return config{addr: addr}, nil
}
