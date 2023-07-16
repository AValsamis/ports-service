package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"ports-service/internal/infra/repository/inmemory"
	"syscall"

	"ports-service/internal/ports/service"
)

func main() {
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)

	repo := inmemory.NewPortRepository()
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

	ctx := context.Background()

	err = srv.LoadPorts(ctx, file, terminateCh)
	if err != nil {
		fmt.Printf("Failed to load ports from file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ports imported: %d\n", repo.GetPortsLength(ctx))
}
