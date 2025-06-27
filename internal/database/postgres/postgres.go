package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type postgres struct {
	conn *pgx.Conn
}

func New(ctx context.Context, connString string) (*postgres, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &postgres{conn: conn}, nil
}
