package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/psxzz/dmsecret-backend/internal/repository/secrets"
)

func (s *service) GetSecretByID(ctx context.Context, secretID uuid.UUID) (string, error) {
	payload, err := s.secrets.GetByID(ctx, secretID)
	if err != nil {
		if errors.Is(err, secrets.ErrNotFound) {
			return "", ErrSecretNotFound
		}

		return "", fmt.Errorf("couldn't get secret: %w", err)
	}

	return payload, nil
}
