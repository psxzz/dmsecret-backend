package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/psxzz/dmsecret-backend/internal/config"
)

type postgres struct {
	conn *pgx.Conn
}

func New(cfg *config.Config) (*postgres, error) {
	conn, err := pgx.Connect(context.Background(), cfg.PGConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &postgres{conn: conn}, nil
}
