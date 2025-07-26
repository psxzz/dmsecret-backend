package secrets

import (
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type Cryptographer interface {
	Encrypt(payload string) (string, error)
	Decrypt(encrypted string) (string, error)
}
type secretsRepository struct {
	kv            valkey.Client
	cryptographer Cryptographer
}

func New(connString string, cryptographer Cryptographer) (*secretsRepository, error) {
	options, err := valkey.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse connection string: %w", err)
	}

	kv, err := valkey.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("couldn't create valkey client: %w", err)
	}

	return &secretsRepository{
		kv:            kv,
		cryptographer: cryptographer,
	}, nil
}

func (r *secretsRepository) Close() {
	r.kv.Close()
}
