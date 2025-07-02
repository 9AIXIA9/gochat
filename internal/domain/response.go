package domain

type SignupResponse struct {
	UserNumber UserNumber `json:"user_number"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateRoomResponse struct {
	RoomNumber  RoomNumber `json:"room_number,string"`
	RoomName    string     `json:"room_name"`
	Owner       UserNumber `json:"owner,string"`
	Description string     `json:"description"`
	MaxUsers    int        `json:"max_users"`
}

var DefaultMessage = NewSuccessMessage()

// middleware || handler response
var (
	TimeoutMessage      = NewDefaultMessage(CodeTimeout)
	UnauthorizedMessage = NewDefaultMessage(CodeUnauthorized)
	InvalidTokenMessage = NewDefaultMessage(CodeInvalidToken)
)

// usecase logic response
var (
	UserExistMessage     = NewDefaultMessage(CodeUserExist)
	UserNotExistMessage  = NewDefaultMessage(CodeUserNotExist)
	WrongPasswordMessage = NewDefaultMessage(CodeWrongPassword)
	RoomExistMessage     = NewDefaultMessage(CodeRoomExist)
	RoomNotExistMessage  = NewDefaultMessage(CodeRoomNotExist)
	RoomIsFullMessage    = NewDefaultMessage(CodeRoomIsFull)
	WrongSecretMessage   = NewDefaultMessage(CodeWrongSecret)
	HasJoinedMessage     = NewDefaultMessage(CodeHasJoined)
	NotJoinedMessage     = NewDefaultMessage(CodeNotJoined)
)
