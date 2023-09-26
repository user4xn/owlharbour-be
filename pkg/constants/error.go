package constants

import "errors"

var (
	BearerTokenHasError = errors.New("bearer token catch error")
	BearerTokenInvalid  = errors.New("invalid token")

	UserNotFound    = errors.New("User Not Found")
	InvalidPassword = errors.New("Invalid Password")

	ErrorGenerateJwt = errors.New("Error Generate JWT")
	EmptyGenerateJwt = errors.New("Empty Generate JWT")

	ErrorLoadLocationTime = errors.New("Error Load Location Time")

	DuplicateStoreUser = errors.New("Duplicate Store Data User")
	ErrorHashPassword  = errors.New("Error Hash Password")

	NotFoundDataUser = errors.New("Not Found Data User")
)
