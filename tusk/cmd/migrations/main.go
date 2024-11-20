package main

import (
	"log"
	"os"
	"tusk/cmd/migrations/runner"
)

func main() {
	if err := runner.RunMigrations(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}
