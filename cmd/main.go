package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ports-service/internal/infra/repository"
	"ports-service/internal/ports/service"
)

func main() {
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)

	repo := repository.NewPortRepository()
	srv := service.NewPortService(repo)

	// Load ports from the PORTS_JSON_PATH file
	filePath := os.Getenv("PORTS_JSON_PATH")
	if filePath == "" {
		fmt.Println("PORTS_JSON_PATH environment variable not set")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	err = srv.LoadPorts(file, terminateCh)
	if err != nil {
		fmt.Printf("Failed to load ports from file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ports imported: %d\n", repo.GetPortsLength())
}
