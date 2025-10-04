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
			err := r.Error()

			if valkey.IsValkeyNil(err) {
				return "", ErrLocked
			}

			if err != nil {
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

		payload, err := r.cryptographer.Decrypt(hash[hashFieldPayload])
		if err != nil {
			return "", fmt.Errorf("couldn't decrypt payload: %w", err)
		}

		return payload, nil
	}

	payload, err := retry.DoWithData(
		getter,
		retry.Context(ctx),
		retry.Attempts(5),
		retry.Delay(10*time.Millisecond),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, ErrLocked)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("couldn't get by id: %w", err)
	}

	return payload, nil
}
