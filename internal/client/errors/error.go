package errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrUserNameInvalid = errors.New("user name is invalid")
)
