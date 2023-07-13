package service_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"ports-service/internal/infra/repository"
	"ports-service/internal/ports/domain"
	"ports-service/internal/ports/service"
)

func TestPortService_LoadPorts(t *testing.T) {
	repo := repository.NewPortRepository()
	portService := service.NewPortService(repo)

	tempFile, err := os.CreateTemp("", "ports.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = os.WriteFile(tempFile.Name(), []byte(samplePorts), 0644)
	assert.NoError(t, err)

	file, err := os.Open(tempFile.Name())
	assert.NoError(t, err)
	defer file.Close()

	err = portService.LoadPorts(file, nil)
	assert.NoError(t, err)

	// Retrieve the loaded ports from the repository
	loadedPort1, err := repo.GetPortByUNLOC("AEJEA")
	assert.NoError(t, err)
	assert.NotNil(t, loadedPort1)
	expectedPort1 := &domain.Port{
		Name:        "Jebel Ali",
		City:        "Jebel Ali",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{55.0272904, 24.9857145},
		Province:    "Dubai",
		Timezone:    "Asia/Dubai",
		UNLOC:       "AEJEA",
		UNLOCs:      []string{"AEJEA"},
		Code:        "52051",
	}
	assert.True(t, comparePorts(loadedPort1, expectedPort1))

	loadedPort2, err := repo.GetPortByUNLOC("AEJED")
	assert.NoError(t, err)
	assert.NotNil(t, loadedPort2)
	expectedPort2 := &domain.Port{
		Name:        "Jebel Dhanna",
		City:        "Jebel Dhanna",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{52.6126027, 24.1915137},
		Province:    "Abu Dhabi",
		Timezone:    "Asia/Dubai",
		UNLOC:       "AEJED",
		UNLOCs:      []string{"AEJED"},
		Code:        "52050",
	}
	assert.True(t, comparePorts(loadedPort2, expectedPort2))
}

func comparePorts(p1, p2 *domain.Port) bool {
	// Convert ports to JSON and compare
	p1JSON, _ := json.Marshal(p1)
	p2JSON, _ := json.Marshal(p2)
	return bytes.Equal(p1JSON, p2JSON)
}
