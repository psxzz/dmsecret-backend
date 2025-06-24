package secrets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type KeyValueClient interface {
	SetEX(ctx context.Context, key string, value string, ttl int) error
}

type secretsRepository struct {
	client KeyValueClient
}

func New(client KeyValueClient) *secretsRepository {
	return &secretsRepository{client: client}
}

func (kv *secretsRepository) CreateSecret(ctx context.Context, secretID uuid.UUID, payload string, ttl int) error {
	secretKey := secretKeyPrefix + secretID.String()

	err := kv.client.SetEX(ctx, secretKey, payload, ttl)
	if err != nil {
		return fmt.Errorf("could not set ex: %w", err)
	}

	return nil
}
