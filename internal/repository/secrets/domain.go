package secrets

import "errors"

const (
	hashFieldPayload   = "payload"
	hashFieldSeenCount = "seenCount"
)

var (
	ErrLocked   = errors.New("secret locked")
	ErrNotFound = errors.New("secret not found")
)
