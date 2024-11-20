package main

//go:generate go run github.com/go-jet/jet/v2/cmd/jet -dsn=postgres://postgres:postgres@localhost:5434/ratt-api?sslmode=disable -path=../internal/ports/database/gen

import (
	"log"
	"os"
	"tusk/cmd/runner"
)

func main() {
	if err := runner.Serve(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}
