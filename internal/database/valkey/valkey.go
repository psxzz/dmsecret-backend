package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"

	"github.com/psxzz/dmsecret-backend/internal/config"
)

type valkeyDB struct {
	client valkey.Client
}

func New(cfg *config.Config) (*valkeyDB, error) {
	options, err := valkey.ParseURL(cfg.ValkeyConnString)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse connection string: %w", err)
	}

	client, err := valkey.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("couldn't create valkey client: %w", err)
	}

	return &valkeyDB{
		client: client,
	}, nil
}

func (v *valkeyDB) SetEX(ctx context.Context, key string, value string, ttlSeconds int) error {
	resp := v.client.Do(ctx,
		v.client.B().Setex().Key(key).Seconds(int64(ttlSeconds)).Value(value).Build(),
	)

	if err := resp.Error(); err != nil {
		return fmt.Errorf("couldn't set key: %w", err)
	}

	return nil
}
