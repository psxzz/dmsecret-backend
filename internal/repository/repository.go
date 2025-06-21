package repository

import (
	"context"

	"github.com/google/uuid"
)

type KeyValueStorage interface {
	CreateSecret(ctx context.Context, secretID uuid.UUID, payload string, TTL int) error
}

type repository struct {
	kvStorage KeyValueStorage
}

func New(kvStorage KeyValueStorage) *repository {
	return &repository{
		kvStorage: kvStorage,
	}
}
