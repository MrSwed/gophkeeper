package errors

import "errors"

var (
	ErrorWrongAuth       = errors.New("wrong auth")
	ErrorSyncNoKey       = errors.New("sync key required")
	ErrorSyncCreatedDate = errors.New("sync with different created date")
	ErrorNoToken         = errors.New("token required")
	ErrorInvalidToken    = errors.New("invalid token")
)
