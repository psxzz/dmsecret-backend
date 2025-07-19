package service

import (
	"context"

	"github.com/google/uuid"
)

type SecretsRepository interface {
	Create(ctx context.Context, secretID uuid.UUID, payload string, TTL int) error
	GetByID(ctx context.Context, id uuid.UUID) (*string, error)
}

type service struct {
	secrets SecretsRepository
}

func New(secrets SecretsRepository) *service {
	return &service{
		secrets: secrets,
	}
}
