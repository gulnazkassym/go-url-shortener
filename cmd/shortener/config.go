package main

import (
	"flag"
	"os"
)

var (
	flagRunAddr string
	flagRunPath string
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagRunPath, "b", "http://localhost:8080/hevfyegruf", "address and port to run server")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envRunPath := os.Getenv("BASE_URL"); envRunPath != "" {
		flagRunPath = envRunPath
	}
}
