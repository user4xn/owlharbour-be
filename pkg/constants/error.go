package constants

import "errors"

var (
	BearerTokenHasError = errors.New("bearer token catch error")
	BearerTokenInvalid  = errors.New("invalid token")
)
