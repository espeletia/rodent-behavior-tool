package main

import (
	"echoes/cmd/runner"
	"log"
	"os"
)

func main() {
	if err := runner.Serve(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}
