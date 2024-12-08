package main

//go:generate go run github.com/go-jet/jet/v2/cmd/jet -dsn=postgres://postgres:postgres@localhost:5434/ratt-api?sslmode=disable -path=../internal/ports/database/gen
//go:generate go run ./openapi/main.go

import (
	"log"
	"os"
	"tusk/cmd/runner"
)

func main() {
	services := os.Args[1:]
	if len(services) == 0 {
		services = []string{
			os.Getenv("TUSK_MODE"),
		}
	}

	if services[0] == "queue" {
		if err := runner.StartQueue(); err != nil {
			log.Println("Error: ", err)
			os.Exit(1)
		}
	} else {
		if err := runner.Serve(); err != nil {
			log.Println("Error: ", err)
			os.Exit(1)
		}
	}
}
