package constant

import "time"

type CtxKey string

const (
	CtxUserID CtxKey = "userID"

	TokenKey = "tokenKey"

	ExpDuration = time.Hour * 24 * 365

	ServerShutdownTimeout = 30
)
