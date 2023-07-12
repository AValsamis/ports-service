package repository_test

import (
	"testing"

	"ports-service/internal/ports/domain"
	"ports-service/internal/ports/repository"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryPortRepository_UpsertPort(t *testing.T) {
	repo := repository.NewPortRepository()

	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}

	err := repo.UpsertPort(port)
	assert.NoError(t, err, "Expected no error")

	result, err := repo.GetPortByUNLOC("TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}

func TestInMemoryPortRepository_GetPortByUNLOC(t *testing.T) {
	repo := repository.NewPortRepository()

	// Retrieve a non-existing port from the repository
	result, err := repo.GetPortByUNLOC("NONEXISTENT")
	assert.NoError(t, err, "Expected no error")
	assert.Nil(t, result, "Expected a nil port")

	port := domain.Port{
		Name:    "Test Port",
		City:    "Test City",
		Country: "Test Country",
		UNLOC:   "TEST",
	}
	err = repo.UpsertPort(port)
	assert.NoError(t, err, "Expected no error")

	result, err = repo.GetPortByUNLOC("TEST")
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil port")
	assert.Equal(t, port.Name, result.Name, "Expected port name to match")
	assert.Equal(t, port.City, result.City, "Expected port city to match")
	assert.Equal(t, port.Country, result.Country, "Expected port country to match")
}
