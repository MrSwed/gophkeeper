package errors

import "errors"

var (
	ErrLoadProfile = errors.New("error get profile")
	ErrDecode      = errors.New("decode error, check passphrase")
)
