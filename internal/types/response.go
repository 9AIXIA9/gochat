package types

import "gochat/internal/domain"

var DefaultResponse = domain.NewSuccessMessage()

// middleware || handler response
var (
	TimeoutResponse      = domain.NewDefaultMessage(domain.CodeTimeout)
	UnauthorizedResponse = domain.NewDefaultMessage(domain.CodeUnauthorized)
	InvalidTokenResponse = domain.NewDefaultMessage(domain.CodeInvalidToken)
)

// usecase logic response
var (
	UserExistResponse     = domain.NewDefaultMessage(domain.CodeUserExist)
	WrongPasswordResponse = domain.NewDefaultMessage(domain.CodeWrongPassword)
	RoomExistResponse     = domain.NewDefaultMessage(domain.CodeRoomExist)
	RoomNotExistResponse  = domain.NewDefaultMessage(domain.CodeRoomNotExist)
	RoomIsFullResponse    = domain.NewDefaultMessage(domain.CodeRoomIsFull)
	WrongSecretResponse   = domain.NewDefaultMessage(domain.CodeWrongSecret)
	HasJoinedResponse     = domain.NewDefaultMessage(domain.CodeHasJoined)
	NotJoinedResponse     = domain.NewDefaultMessage(domain.CodeNotJoined)
)
