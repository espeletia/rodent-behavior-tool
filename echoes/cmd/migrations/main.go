package main

import (
	"echoes/cmd/migrations/runner"
	"log"
	"os"
)

func main() {
	if err := runner.RunMigrations(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}
