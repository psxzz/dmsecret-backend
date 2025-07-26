package secrets

import (
	"fmt"

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
