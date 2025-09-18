package types

import (
	"context"
	"errors"
)

// repository
var (
	ErrNotFound     = errors.New("entity not found")
	ErrDuplicateKey = errors.New("entity key has existed")
)

//handler

var (
	ErrNullPointer = errors.New("param is nil")
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

// websocket
var (
	ErrClientClosed = errors.New("client has been closed")
	ErrFullMessage  = errors.New("too many messages")
)
