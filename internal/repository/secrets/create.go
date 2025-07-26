package secrets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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
