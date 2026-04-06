package error

import "errors"

var (
	ErrInternal         = errors.New("internal server error")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrUserNotFound     = errors.New("user not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrContextTimeout   = errors.New("request time out")
	ErrRequestCancelled = errors.New("request cancelled")
	ErrBadRequest       = errors.New("bad request")
)
