package secrets

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
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

	getter := func() (string, error) {
		resp := client.Do(ctx, client.B().Watch().Key(key).Build())
		if err := resp.Error(); err != nil {
			return "", fmt.Errorf("couldn't watch key: %w", err)
		}

		resp = client.Do(ctx, client.B().Exists().Key(key).Build())
		if err := resp.Error(); err != nil {
			return "", fmt.Errorf("couldn't check key existence: %w", err)
		}

		exists, err := resp.AsBool()
		if err != nil {
			return "", fmt.Errorf("couldn't get bool from response: %w", err)
		}

		if !exists {
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
				return "", fmt.Errorf("couldn't do multi pipeline: %w", err)
			}
		}

		resp = client.Do(ctx, client.B().Hgetall().Key(key).Build())
		if err := resp.Error(); err != nil {
			return "", fmt.Errorf("couldn't get key: %w", err)
		}

		hash, err := resp.AsStrMap()
		if err != nil {
			return "", fmt.Errorf("couldn't get hash from response: %w", err)
		}

		seenCount, err := strconv.Atoi(hash[hashFieldSeenCount])
		if err != nil {
			return "", fmt.Errorf("couldn't convert seen counter: %w", err)
		}

		if seenCount == 0 {
			resp := client.Do(ctx, client.B().Unlink().Key(key).Build())
			if err := resp.Error(); err != nil {
				return "", fmt.Errorf("couldn't unlink: %w", err)
			}
		}

		return hash[hashFieldPayload], nil
	}

	payload, err := retry.DoWithData(
		getter,
		retry.Context(ctx),
		retry.Attempts(5),
		retry.Delay(10*time.Millisecond),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, ErrNotFound)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("couldn't get by id: %w", err)
	}

	return payload, nil
}
