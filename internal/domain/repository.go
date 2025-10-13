package domain

type Repository interface {
	Users() UserRepository
	Rooms() RoomRepository
	Messages() MessageRepository
	RefreshTokens() RefreshTokenRepository
}

type UserRepository interface {
	UserSaver
	UserFinder
}

type RoomRepository interface {
	RoomSaver
	RoomFinder
	RoomJoiner
	RoomLeaver
	RoomMemberFinder
	JoinRoomAggregateUOW
}

type MessageRepository interface {
	MessageSaver
	UnsentMessageFinder
	SentMessageUpdater
}

type RefreshTokenRepository interface {
	RefreshTokenSaver
	RefreshTokenFinder
}
