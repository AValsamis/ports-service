package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"ports-service/internal/infra/repository/redis"
	"syscall"

	"ports-service/internal/ports/service"
)

func main() {
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Println("REDIS_URL environment variable not set")
		os.Exit(1)
	}

	repo, err := redis.NewPortRepository(redisURL)
	if err != nil {
		fmt.Printf("Failed to create Redis repository: %v\n", err)
		os.Exit(1)
	}
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

	length, err := repo.GetPortsLength(ctx)
	if err != nil {
		fmt.Printf("Failed to get ports length: %v\n", err)
	}
	fmt.Printf("Ports imported: %d\n", length)
}
