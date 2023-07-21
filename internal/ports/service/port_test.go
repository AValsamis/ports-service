package service_test

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"ports-service/internal/ports/domain"
	"ports-service/internal/ports/service"
)

type mockPortRepository struct {
	ports map[string]domain.Port
}

func (m *mockPortRepository) GetPortByUNLOC(_ context.Context, unloc string) (*domain.Port, error) {
	port, exists := m.ports[unloc]
	if !exists {
		return nil, nil
	}
	return &port, nil
}

func (m *mockPortRepository) UpsertPort(_ context.Context, port domain.Port) error {
	m.ports[port.UNLOC] = port
	return nil
}

func TestPortService_LoadPorts_Success(t *testing.T) {
	ctx := context.Background()

	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	portService := service.NewPortService(repo)

	reader := bytes.NewReader([]byte(samplePorts))

	err := portService.LoadPorts(ctx, reader, nil)
	assert.NoError(t, err)

	portAEJEA, err := repo.GetPortByUNLOC(ctx, "AEJEA")
	assert.NoError(t, err)
	assert.NotNil(t, portAEJEA)
	assert.Equal(t, "Jebel Ali", portAEJEA.Name)
	assert.Equal(t, "Jebel Ali", portAEJEA.City)
	assert.Equal(t, "United Arab Emirates", portAEJEA.Country)

	portAEJED, err := repo.GetPortByUNLOC(ctx, "AEJED")
	assert.NoError(t, err)
	assert.NotNil(t, portAEJED)
	assert.Equal(t, "Jebel Dhanna", portAEJED.Name)
	assert.Equal(t, "Jebel Dhanna", portAEJED.City)
	assert.Equal(t, "United Arab Emirates", portAEJED.Country)
}

func TestPortService_LoadPorts_Error(t *testing.T) {
	ctx := context.Background()

	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	portService := service.NewPortService(repo)

	inputWithInvalidPort := samplePorts + `
		"InvalidPort": {
			"name": "Invalid Port"
		}
	}`

	reader := bytes.NewReader([]byte(inputWithInvalidPort))
	err := portService.LoadPorts(ctx, reader, nil)
	assert.NoError(t, err)

	// Verify that the ports without validation errors are still upserted
	portAEJEA, err := repo.GetPortByUNLOC(ctx, "AEJEA")
	assert.NoError(t, err)
	assert.NotNil(t, portAEJEA)
	assert.Equal(t, "Jebel Ali", portAEJEA.Name)
	assert.Equal(t, "Jebel Ali", portAEJEA.City)
	assert.Equal(t, "United Arab Emirates", portAEJEA.Country)

	portAEJED, err := repo.GetPortByUNLOC(ctx, "AEJED")
	assert.NoError(t, err)
	assert.NotNil(t, portAEJED)
	assert.Equal(t, "Jebel Dhanna", portAEJED.Name)
	assert.Equal(t, "Jebel Dhanna", portAEJED.City)
	assert.Equal(t, "United Arab Emirates", portAEJED.Country)

	// Verify that the port with validation errors is not upserted
	invalidPort, err := repo.GetPortByUNLOC(ctx, "InvalidPort")
	assert.NoError(t, err)
	assert.Nil(t, invalidPort)
}

func TestPortService_LoadPorts_Graceful_Termination(t *testing.T) {
	ctx := context.Background()

	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	portService := service.NewPortService(repo)

	reader := bytes.NewReader([]byte(samplePorts))

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	// Send a termination signal to terminate the method prematurely
	terminate <- syscall.SIGTERM

	err := portService.LoadPorts(ctx, reader, terminate)
	assert.NoError(t, err)

	// Verify that no ports are upserted due to premature termination
	portAEJEA, err := repo.GetPortByUNLOC(ctx, "AEJEA")
	assert.NoError(t, err)
	assert.Nil(t, portAEJEA)

	portAEJED, err := repo.GetPortByUNLOC(ctx, "AEJED")
	assert.NoError(t, err)
	assert.Nil(t, portAEJED)
}

var samplePorts = `{
		"AEJEA": {
			"name": "Jebel Ali",
			"city": "Jebel Ali",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": [55.0272904, 24.9857145],
			"province": "Dubai",
			"timezone": "Asia/Dubai",
			"unlocs": ["AEJEA"],
			"code": "52051"
		},
		"AEJED": {
			"name": "Jebel Dhanna",
			"city": "Jebel Dhanna",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"coordinates": [52.6126027, 24.1915137],
			"province": "Abu Dhabi",
			"timezone": "Asia/Dubai",
			"unlocs": ["AEJED"],
			"code": "52050"
		}
	}`
