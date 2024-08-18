package errors

import "errors"

var (
	ErrLoadProfile     = errors.New("error get profile")
	ErrDecode          = errors.New("decode error, check passphrase")
	ErrPassword        = errors.New("wrong password")
	ErrPasswordConfirm = errors.New("password confirm error")
)
