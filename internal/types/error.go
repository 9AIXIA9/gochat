package types

import (
	"context"
	"errors"
)

//handler

var (
	ErrNullPointer = errors.New("param is nil")
)

// websocket
var (
	ErrUserExist    = errors.New("the user has existed")
	ErrRoomExist    = errors.New("the room has existed")
	ErrRoomNotExist = errors.New("the room doesn't exist")
)

// jwt
var (
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrInvalidToken       = errors.New("invalid token")
)

// timeout
var (
	ErrTimeout  = context.DeadlineExceeded
	ErrCanceled = context.Canceled
)
