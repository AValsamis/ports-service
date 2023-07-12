package repository

import (
	"sync"

	"ports-service/internal/ports/domain"
)

// InMemoryPortRepository is an in-memory repository handling ports.
type InMemoryPortRepository struct {
	ports map[string]domain.Port
	mutex sync.RWMutex
}

// NewPortRepository creates a new instance of InMemoryPortRepository.
func NewPortRepository() *InMemoryPortRepository {
	return &InMemoryPortRepository{
		ports: make(map[string]domain.Port),
	}
}

// UpsertPort inserts or updates a port in the repository.
func (r *InMemoryPortRepository) UpsertPort(port domain.Port) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.ports[port.UNLOC] = port
	return nil
}

// GetPortByUNLOC retrieves a port from the repository by its UNLOC code.
func (r *InMemoryPortRepository) GetPortByUNLOC(unloc string) (*domain.Port, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	port, exists := r.ports[unloc]
	if !exists {
		return nil, nil
	}
	return &port, nil
}

// GetPortsLength returns the total number of ports in the repository.
func (r *InMemoryPortRepository) GetPortsLength() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.ports)
}
