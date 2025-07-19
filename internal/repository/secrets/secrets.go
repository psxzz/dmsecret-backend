package secrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/psxzz/dmsecret-backend/internal/database/valkey"
)

type KeyValueClient interface {
	Get(ctx context.Context, key string) (string, error)
	SetEX(ctx context.Context, key string, value string, ttl int) error
}

type secretsRepository struct {
	kv KeyValueClient
}

func New(kv KeyValueClient) *secretsRepository {
	return &secretsRepository{kv: kv}
}

func (r *secretsRepository) Create(ctx context.Context, secretID uuid.UUID, payload string, ttl int) error {
	secretKey := getSecretKey(secretID)

	err := r.kv.SetEX(ctx, secretKey, payload, ttl)

	if err != nil {
		return fmt.Errorf("could not set ex: %w", err)
	}

	return nil
}

func (r *secretsRepository) GetByID(ctx context.Context, secretID uuid.UUID) (*string, error) {
	secretKey := getSecretKey(secretID)

	payload, err := r.kv.Get(ctx, secretKey)
	if err != nil {
		if errors.Is(err, valkey.ErrNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not get: %w", err)
	}

	return &payload, nil
}
