package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"ports-service/internal/ports/domain"
)

// PortRepository is an interface for accessing port data.
type PortRepository interface {
	GetPortByUNLOC(unloc string) (*domain.Port, error)
	UpsertPort(port domain.Port) error
}

// PortService provides methods for managing ports.
type PortService struct {
	repo PortRepository
}

// NewPortService creates a new instance of PortService.
func NewPortService(repo PortRepository) *PortService {
	return &PortService{
		repo: repo,
	}
}

// LoadPorts loads ports from the given reader.
// It expects the input to be in JSON format.
// Each JSON object should represent a single port.
// It decodes the JSON data from the reader and upserts the ports into the repository.
func (s *PortService) LoadPorts(reader io.Reader, terminate chan os.Signal) error {
	dec := json.NewDecoder(reader)

	t, err := dec.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return fmt.Errorf("expected {, got %v", t)
	}

	for dec.More() {
		// Check for termination signal
		select {
		case <-terminate:
			return nil // Gracefully terminate
		default:
			// Continue processing
		}

		t, err := dec.Token()
		if err != nil {
			return err
		}
		key := t.(string)

		var port domain.Port
		if err := dec.Decode(&port); err != nil {
			return err
		}

		port.UNLOC = key
		err = s.upsertPort(port)
		if err != nil {
			log.Warnf("failed to upsert port: %v", err)
			continue
		}
	}
	return nil
}

// upsertPort inserts or updates a port in the repository.
// It validates the port's data before upserting.
// If the validation fails, it returns an error.
func (s *PortService) upsertPort(port domain.Port) error {
	if err := port.Validate(); err != nil {
		return err
	}

	return s.repo.UpsertPort(port)
}
