package secrets

import "errors"

const (
	hashFieldPayload   = "payload"
	hashFieldSeenCount = "seenCount"
)

var ErrNotFound = errors.New("secret not found")
