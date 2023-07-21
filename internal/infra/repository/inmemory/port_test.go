package inmemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"ports-service/internal/infra/repository/inmemory"
	"ports-service/internal/ports/domain"
)

func TestInMemoryPortRepository_UpsertPort(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewPortRepository()

	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}

	err := repo.UpsertPort(ctx, port)
	assert.NoError(t, err, "Expected no error")

	result, err := repo.GetPortByUNLOC(ctx, "TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}

func TestInMemoryPortRepository_GetPortByUNLOC(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewPortRepository()

	// Retrieve a non-existing port from the repository
	result, err := repo.GetPortByUNLOC(ctx, "NONEXISTENT")
	assert.NoError(t, err, "Expected no error")
	assert.Nil(t, result, "Expected a nil port")

	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}
	err = repo.UpsertPort(ctx, port)
	assert.NoError(t, err, "Expected no error")

	result, err = repo.GetPortByUNLOC(ctx, "TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}
