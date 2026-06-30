package main

import (
	"fmt"
	"os"

	"github.com/faiyaz032/aerokv/internal/server"
)

const (
	defaultHost = "localhost"
	defaultPort = "9111"
)

func main() {
	srv := server.New(defaultHost, defaultPort)

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server lifecycle error: %v\n", err)
		os.Exit(1)
	}
}
