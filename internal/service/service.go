package service

import (
	"context"

	"github.com/google/uuid"
)

type SecretsRepository interface {
	CreateSecret(ctx context.Context, secretID uuid.UUID, payload string, TTL int) error
}

type service struct {
	secrets SecretsRepository
}

func New(secrets SecretsRepository) *service {
	return &service{
		secrets: secrets,
	}
}
