package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const defaultSecretTTL = 3600

func (s *service) CreateSecret(ctx context.Context, payload string) (string, error) {
	secretID := uuid.New()

	err := s.secrets.Create(ctx, secretID, payload, defaultSecretTTL)
	if err != nil {
		return "", fmt.Errorf("couldn't create secret: %w", err)
	}

	return secretID.String(), nil
}
