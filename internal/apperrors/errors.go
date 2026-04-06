package apperrors

import "errors"

var (
	ErrInternal           = errors.New("internal server error")
	ErrUserAlreadyExists  = errors.New("user already exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrDeadlineExceeded   = errors.New("request time out")
	ErrCanceled           = errors.New("request cancelled")
	ErrBadRequest         = errors.New("bad request")
)
