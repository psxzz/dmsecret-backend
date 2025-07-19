package valkey

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/valkey-io/valkey-go"
)

var ErrNotFound = errors.New("not found")

type valkeyDB struct {
	client valkey.Client
}

func New(connString string) (*valkeyDB, error) {
	options, err := valkey.ParseURL(connString)
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

func (v *valkeyDB) SetEX(ctx context.Context, key string, payload string, ttlSeconds int) error {
	resps := v.client.DoMulti(ctx,
		v.client.B().Hset().Key(key).FieldValue().
			FieldValue("payload", payload).
			FieldValue("seen_count", "1").Build(),
		v.client.B().Expire().Key(key).Seconds(int64(ttlSeconds)).Build(),
	)

	for _, resp := range resps {
		if err := resp.Error(); err != nil {
			return fmt.Errorf("couldn't set key: %w", err)
		}
	}

	return nil
}

func (v *valkeyDB) Get(ctx context.Context, key string) (string, error) {
	client, cancel := v.client.Dedicate()
	defer cancel()

	resp := client.Do(ctx, client.B().Watch().Key(key).Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("couldn't watch key: %w", err)
	}

	resp = client.Do(ctx, client.B().Hgetall().Key(key).Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("couldn't get key: %w", err)
	}

	hash, err := resp.AsStrMap()
	if err != nil {
		return "", fmt.Errorf("couldn't get hash from response: %w", err)
	}

	if len(hash) == 0 {
		return "", ErrNotFound
	}

	resps := client.DoMulti(
		ctx,
		client.B().Multi().Build(),
		client.B().Hincrby().Key(key).Field("seen_count").Increment(-1).Build(),
		client.B().Exec().Build(),
	)
	for _, r := range resps {
		if err := r.Error(); err != nil {
			return "", fmt.Errorf("couldn't do pipelined: %w", err)
		}
	}

	resp = client.Do(ctx, client.B().Hgetall().Key(key).Build())
	if err := resp.Error(); err != nil {
		return "", fmt.Errorf("couldn't get key: %w", err)
	}

	hash, err = resp.AsStrMap()
	if err != nil {
		return "", fmt.Errorf("couldn't get hash from response: %w", err)
	}

	seenCount, err := strconv.Atoi(hash["seen_count"])
	if err != nil {
		return "", fmt.Errorf("couldn't convert seen counter: %w", err)
	}

	if seenCount < 0 {
		client.Do(ctx, client.B().Unlink().Key(key).Build())

		return "", ErrNotFound
	}

	return hash["payload"], nil
}
