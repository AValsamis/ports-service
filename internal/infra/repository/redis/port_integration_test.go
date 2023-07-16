package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"ports-service/internal/infra/repository/redis"
	"ports-service/internal/ports/domain"
)

var (
	redisContainer testcontainers.Container
	redisRepo      *redis.PortRepository
)

func setupRedisContainer(t *testing.T) (func(), error) {
	ctx := context.Background()

	// Define the Redis container configuration
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	// Create and start the Redis container
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Redis container: %w", err)
	}

	// Get the Redis container's host and port
	redisHost, err := redisC.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis container host: %w", err)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis container port: %w", err)
	}

	// Create the RedisPortRepository with the container's host and port
	redisURL := fmt.Sprintf("redis://%s:%s/0", redisHost, redisPort.Port())
	repo, err := redis.NewPortRepository(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis repository: %w", err)
	}

	// Set the global variables for the container and repository
	redisContainer = redisC
	redisRepo = repo

	// Create a cleanup function to terminate the container after the tests
	cleanup := func() {
		err := redisContainer.Terminate(ctx)
		if err != nil {
			t.Errorf("failed to terminate Redis container: %v", err)
		}
	}

	return cleanup, nil
}

func TestRedisPortRepository_UpsertPort(t *testing.T) {
	cleanup, err := setupRedisContainer(t)
	assert.NoError(t, err, "Failed to set up Redis container")
	defer cleanup()

	ctx := context.Background()
	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}

	err = redisRepo.UpsertPort(ctx, port)
	assert.NoError(t, err, "Expected no error")

	result, err := redisRepo.GetPortByUNLOC(ctx, "TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}

func TestRedisPortRepository_GetPortByUNLOC(t *testing.T) {
	ctx := context.Background()
	cleanup, err := setupRedisContainer(t)
	assert.NoError(t, err, "Failed to set up Redis container")
	defer cleanup()

	// Retrieve a non-existing port from the repository
	result, err := redisRepo.GetPortByUNLOC(ctx, "NONEXISTENT")
	assert.NoError(t, err, "Expected no error")
	assert.Nil(t, result, "Expected a nil port")

	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}
	err = redisRepo.UpsertPort(ctx, port)
	assert.NoError(t, err, "Expected no error")

	result, err = redisRepo.GetPortByUNLOC(ctx, "TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}
