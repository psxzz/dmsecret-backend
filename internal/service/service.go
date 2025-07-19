package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrSecretNotFound = errors.New("secret not found")

type SecretsRepository interface {
	Create(ctx context.Context, secretID uuid.UUID, payload string, TTL int) error
	GetByID(ctx context.Context, id uuid.UUID) (string, error)
}

type service struct {
	secrets SecretsRepository
}

func New(secrets SecretsRepository) *service {
	return &service{
		secrets: secrets,
	}
}
