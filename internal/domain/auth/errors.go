package auth

import "errors"

var (
	ErrPasswordVerifyFail      = errors.New("Password verification failed")
	ErrTokenInvalid            = errors.New("Token is invalid")
	ErrTokenExpired            = errors.New("Token is expired")
	ErrTokenMissing            = errors.New("Token is missing")
	ErrTokenMalformed          = errors.New("Token is malformed")
	ErrTokenNotActive          = errors.New("Token is not active yet")
	ErrTokenWrongSigningMethod = errors.New("Token has wrong signing method")
	ErrTokenGenerateFail       = errors.New("Failed to generate token")
	ErrorRefreshTokenEmpty     = errors.New("Refresh Token empty string")
)
