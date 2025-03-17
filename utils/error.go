package utils

import "errors"

var (
	ErrFailedRegister = errors.New("failed register")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrRoomNotFound   = errors.New("room doesn't exist")
	ErrBadRequest     = errors.New("bad request")
	ErrInternal       = errors.New("internal server error")
	ErrUserNotFound   = errors.New("internal server error")
)
