package errors

import "errors"

// Repository errors
var (
	ErrCartNotFound = errors.New("cart not found")
	ErrItemNotFound = errors.New("item not found")
	ErrUserNotFound = errors.New("user not found")
)

// Service errors
var (
	ErrUserExists         = errors.New("user exists")
	ErrUserCreationFailed = errors.New("user creation failed")
	ErrUserVerifyFailed   = errors.New("user verification failed")
	ErrInvalidCreds       = errors.New("invalid credentials")

	ErrTokenGeneration     = errors.New("token generation failed")
	ErrTokenStorage        = errors.New("token storage failed")
	ErrTokenParsingFailed  = errors.New("token parsing failed")
	ErrInvalidTokenSub     = errors.New("invalid token subject")
	ErrTokenDeletionFailed = errors.New("token deletion failed")
	ErrTokenSaveFailed     = errors.New("token save failed")
)

// Handler errors and messages
var (
	ErrInvalidID = errors.New("invalid id")
)

// Middleware errors
var (
	ErrClaimsNotFound    = errors.New("claims not found")
	ErrInvalidClaims     = errors.New("invalid claims")
	ErrAdminOnly         = errors.New("admin only")
	ErrAuthHeaderMissing = errors.New("authorization header missing or invalid")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenTypeFailed   = errors.New("token type assertion failed")
	ErrUnauthorizedToken = errors.New("unauthorized or invalid token")
)
