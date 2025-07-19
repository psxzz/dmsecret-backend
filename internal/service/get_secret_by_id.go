package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *service) GetSecretByID(ctx context.Context, secretID uuid.UUID) (*string, error) {
	payload, err := s.secrets.GetByID(ctx, secretID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get secret: %w", err)
	}

	return payload, nil
}
