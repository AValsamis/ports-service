package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"ports-service/internal/ports/domain"
)

// PortRepository is a Redis repository handling ports.
type PortRepository struct {
	client *redis.Client
}

// NewPortRepository creates a new instance of PortRepository.
func NewPortRepository(redisURL string) (*PortRepository, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	return &PortRepository{
		client: client,
	}, nil
}

// UpsertPort inserts or updates a port in the repository.
func (r *PortRepository) UpsertPort(ctx context.Context, port domain.Port) error {
	data, err := json.Marshal(port)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, port.UNLOC, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetPortByUNLOC retrieves a port from the repository by its UNLOC code.
func (r *PortRepository) GetPortByUNLOC(ctx context.Context, unloc string) (*domain.Port, error) {
	data, err := r.client.Get(ctx, unloc).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var port domain.Port
	err = json.Unmarshal(data, &port)
	if err != nil {
		return nil, err
	}

	return &port, nil
}

// GetPortsLength returns the total number of ports in the repository.
func (r *PortRepository) GetPortsLength(ctx context.Context) (int64, error) {
	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return 0, err
	}

	return int64(len(keys)), nil
}
