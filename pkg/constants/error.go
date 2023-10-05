package constants

import "errors"

var (
	BearerTokenHasError = errors.New("Bearer token catch error")
	BearerTokenInvalid  = errors.New("Invalid token")

	UserNotFound    = errors.New("User not found")
	InvalidPassword = errors.New("Invalid password")

	ErrorGenerateJwt = errors.New("Error generate JWT")
	EmptyGenerateJwt = errors.New("Empty generate JWT")

	ErrorLoadLocationTime = errors.New("Error load location time")

	DuplicateStoreUser = errors.New("Duplicate store data user")
	ErrorHashPassword  = errors.New("Error hash password")

	NotFoundDataUser = errors.New("Not found data user")
	FailedUpdateUser = errors.New("Failed update user")
	FailedDeleteUser = errors.New("Failed delete user")

	FailedChangePassword   = errors.New("Failed change password")
	FailedNotSamePassword  = errors.New("Please confirm the same password")
	MinimCharacterPassword = errors.New("Minimum password is 8 characters")
	PasswordSameCurrent    = errors.New("The password is the same as the current one")
	ErrorDecodeBase64      = errors.New("Sorry failed to decode base64")
	FailedVerifyEmail      = errors.New("Sorry failed to verify email")
	UserNotVerifyEmail     = errors.New("Please verification account simpel!")
)
