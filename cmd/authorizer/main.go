package main

import (
	"authorizer/internal/app/service"
	"authorizer/internal/app/storage"
	"os"

	log "github.com/sirupsen/logrus"

	cmd "authorizer/cmd/root"
	"authorizer/internal/common/logfile"
)

func main() {
	logfile.Init()

	log.Println("--- starting app ---")

	// Initialize DB
	db := storage.InMemory{}

	// Initialize service
	svc := service.New(&db)

	// Get input from stdin
	stdin := os.Stdin
	// return output to stdout
	stdout := os.Stdout

	// Execute application
	cmd.Execute(svc, stdin, stdout)
}
