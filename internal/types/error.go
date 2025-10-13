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
	ErrFullMessage = errors.New("messages is full")
)

var (
	ErrLengthLessThanZero = errors.New("length should be greater than 0")
	ErrEmptyPointer       = errors.New("empty pointer")
)
