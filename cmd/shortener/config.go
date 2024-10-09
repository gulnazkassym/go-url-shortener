package main

import (
	"flag"
)

var (
	flagRunAddr string
	flagRunPath string
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagRunPath, "b", "http://localhost:8000/hevfyegruf", "address and port to run server")

	flag.Parse()
}
