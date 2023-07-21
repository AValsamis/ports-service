package inmemory

import (
	"context"
	"sync"

	"ports-service/internal/ports/domain"
)

// PortRepository is an in-memory repository handling ports.
type PortRepository struct {
	ports map[string]domain.Port
	mutex sync.RWMutex
}

// NewPortRepository creates a new instance of InMemoryPortRepository.
func NewPortRepository() *PortRepository {
	return &PortRepository{
		ports: make(map[string]domain.Port),
	}
}

// UpsertPort inserts or updates a port in the repository.
func (r *PortRepository) UpsertPort(_ context.Context, port domain.Port) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.ports[port.UNLOC] = port
	return nil
}

// GetPortByUNLOC retrieves a port from the repository by its UNLOC code.
func (r *PortRepository) GetPortByUNLOC(_ context.Context, unloc string) (*domain.Port, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	port, exists := r.ports[unloc]
	if !exists {
		return nil, nil
	}
	return &port, nil
}

// GetPortsLength returns the total number of ports in the repository.
func (r *PortRepository) GetPortsLength(ctx context.Context) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.ports)
}
