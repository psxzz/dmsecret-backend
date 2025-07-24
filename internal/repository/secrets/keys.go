package secrets

import "github.com/google/uuid"

const secretKeyPrefix = "sec1:"

func getSecretKey(id uuid.UUID) string {
	return secretKeyPrefix + id.String()
}
