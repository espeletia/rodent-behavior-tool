package main

//go:generate templ generate -path ./..

import (
	"log"
	"os"
	"valentine/cmd/runner"
)

func main() {
	err := runner.Serve()
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}
