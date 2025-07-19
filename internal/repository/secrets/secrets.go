package secrets

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type secretsRepository struct {
	kv valkey.Client
}

func New(connString string) (*secretsRepository, error) {
	options, err := valkey.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse connection string: %w", err)
	}

	kv, err := valkey.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("couldn't create valkey client: %w", err)
	}

	return &secretsRepository{kv: kv}, nil
}

func (r *secretsRepository) Create(ctx context.Context, secretID uuid.UUID, payload string, ttl int) error {
	key := getSecretKey(secretID)

	resps := r.kv.DoMulti(
		ctx,
		r.kv.B().Hset().Key(key).FieldValue().
			FieldValue(hashFieldPayload, payload).
			FieldValue(hashFieldSeenCount, "1").Build(),
		r.kv.B().Expire().Key(key).Seconds(int64(ttl)).Build(),
	)
	for _, resp := range resps {
		if err := resp.Error(); err != nil {
			return fmt.Errorf("couldn't set key: %w", err)
		}
	}

	return nil
}

func (r *secretsRepository) GetByID(ctx context.Context, secretID uuid.UUID) (string, error) {
	client, cancel := r.kv.Dedicate()
	defer cancel()

	key := getSecretKey(secretID)

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
		client.B().Hincrby().Key(key).Field(hashFieldSeenCount).Increment(-1).Build(),
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

	seenCount, err := strconv.Atoi(hash[hashFieldSeenCount])
	if err != nil {
		return "", fmt.Errorf("couldn't convert seen counter: %w", err)
	}

	if seenCount < 0 {
		client.Do(ctx, client.B().Unlink().Key(key).Build())

		return "", ErrNotFound
	}

	return hash[hashFieldPayload], nil
}
