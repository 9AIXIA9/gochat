package di

import (
	rootapp "gochat/internal/application"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	socialApp "gochat/internal/social/application"
	socialDomain "gochat/internal/social/domain"
	socialPersistence "gochat/internal/social/infrastructure/persistence/repository"

	"github.com/google/wire"
)

var UseCaseHTTPSet = wire.NewSet(
	provideSignUpUseCase,
	provideLoginUseCase,
	provideRefreshAccessTokenUseCase,
	provideParseAccessTokenUseCase,
	provideSendPrivateMessageUseCase,
	provideSendRoomMessageUseCase,
	provideCreateRoomUseCase,
	provideJoinRoomUseCase,
	provideLeaveRoomUseCase,
)

var UseCaseWebsocketSet = wire.NewSet(
	provideWebsocketUserSessionStartedUseCase,
	provideNotificationPrivateMessageReadUseCase,
	provideNotificationRoomMessageReadUseCase,
)

var UseCaseKafkaSet = wire.NewSet(
	provideAuthUserCreatedUseCase,
	provideSocialUserCreatedUseCase,
	provideSocialRoomCreatedUseCase,
	provideSocialRoomJoinedUseCase,
	provideSocialRoomLeftUseCase,
	provideChatUserCreatedUseCase,
	provideChatRoomCreatedUseCase,
	provideChatRoomJoinedUseCase,
	provideChatRoomLeftUseCase,
	provideChatRoomMessageCreatedUseCase,
	provideChatPrivateMessageCreatedUseCase,
	provideNotificationWelcomeEmailNotificationRequestedUseCase,
	provideNotificationRoomMessageNotificationRequestedUseCase,
	provideNotificationPrivateMessageNotificationRequestedUseCase,
	provideNotificationUndeliveredMessagesNotificationRequestedUseCase,
)

var UseCaseBinlogReaderSet = wire.NewSet(
	provideUnpublishedEventsCreatedCase,
)

// -------------------- UseCases (HTTP side) --------------------
func provideUnpublishedEventsCreatedCase(
	publisher event.Publisher,
	eventRepo event.Repository,
) rootapp.UnpublishedEventsCreatedUseCase {
	return rootapp.NewUnpublishedEventsCreatedUseCase(publisher, eventRepo)
}

// -------------------- UseCases (HTTP side) --------------------
func provideSignUpUseCase(
	eventIDGen event.IDGenerator,
	userIDGen authDomain.UserIDGenerator,
	numberGen authDomain.UserNumberGenerator,
	encryptor authDomain.Encryptor,
	userRepo authDomain.UserRepository,
) authApp.SignUpUseCase {
	return authApp.NewSignUpUseCase(eventIDGen, userIDGen, numberGen, encryptor, userRepo)
}

func provideLoginUseCase(
	eventIDGen event.IDGenerator,
	comparator authDomain.Comparator,
	accessTokenGenerator authDomain.AccessTokenGenerator,
	refreshTokenGenerator authDomain.RefreshTokenGenerator,
	userRepo authDomain.UserRepository,
	refreshTokenRepo authDomain.RefreshTokenRepository,
) authApp.LoginUseCase {
	return authApp.NewLoginUseCase(eventIDGen, comparator, userRepo, refreshTokenRepo, accessTokenGenerator, refreshTokenGenerator)
}
func provideRefreshAccessTokenUseCase(
	refreshTokenRepo authDomain.RefreshTokenRepository,
	accessTokenGenerator authDomain.AccessTokenGenerator,
	refreshTokenGenerator authDomain.RefreshTokenGenerator,
) authApp.RefreshAccessTokenUseCase {
	return authApp.NewRefreshAccessTokenUseCase(refreshTokenRepo, refreshTokenRepo, accessTokenGenerator, refreshTokenGenerator)
}
func provideParseAccessTokenUseCase(accessTokenParser authApp.AccessTokenParser) authApp.ParseAccessTokenUseCase {
	return authApp.NewParseAccessTokenUseCase(accessTokenParser)
}
func provideSendPrivateMessageUseCase(
	messageIDGen chatDomain.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	userRepo chatDomain.UserRepository,
	messageRepo chatDomain.PrivateMessageRepository,
) chatApp.SendPrivateMessageUseCase {
	return chatApp.NewSendPrivateMessageUseCase(messageIDGen, eventIDGen, userRepo, messageRepo)
}
func provideSendRoomMessageUseCase(
	messageIDGen chatDomain.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomRepo chatDomain.RoomRepository,
	messageRepo chatDomain.RoomMessageRepository,
) chatApp.SendRoomMessageUseCase {
	return chatApp.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomRepo, messageRepo)
}
func provideCreateRoomUseCase(
	eventIDGen event.IDGenerator,
	roomIDGen socialDomain.RoomIDGenerator,
	numberGen socialDomain.RoomNumberGenerator,
	encryptor socialDomain.Encryptor,
	roomRepo socialDomain.RoomRepository,
) socialApp.CreateRoomUseCase {
	return socialApp.NewCreateRoomUseCase(eventIDGen, roomIDGen, numberGen, encryptor, roomRepo)
}
func provideJoinRoomUseCase(
	eventIDGen event.IDGenerator,
	comparator socialDomain.Comparator,
	roomRepo socialDomain.RoomRepository,
) socialApp.JoinRoomUseCase {
	return socialApp.NewJoinRoomUseCase(eventIDGen, roomRepo, comparator, roomRepo)
}
func provideLeaveRoomUseCase(
	eventIDGen event.IDGenerator,
	roomRepo socialDomain.RoomRepository,
) socialApp.LeaveRoomUseCase {
	return socialApp.NewLeaveRoomUseCase(eventIDGen, roomRepo, roomRepo)
}

// -------------------- Event UseCases (websocket side) --------------------
func provideWebsocketUserSessionStartedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) rootapp.UserSessionStartedUseCase {
	return rootapp.NewUserSessionStartedUseCase(eventIDGen, eventRepo)
}
func provideNotificationPrivateMessageReadUseCase(messageRepo notificationDomain.PrivateMessageRepository) notificationApp.PrivateMessageReadUseCase {
	return notificationApp.NewPrivateMessageReadUseCase(messageRepo, messageRepo)
}
func provideNotificationRoomMessageReadUseCase(messageRepo notificationDomain.RoomMessageRepository) notificationApp.RoomMessageReadUseCase {
	return notificationApp.NewRoomMessageReadUseCase(messageRepo, messageRepo)
}

// -------------------- Event UseCases (Kafka consumer side) --------------------

func provideAuthUserCreatedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
	userRepo authDomain.UserRepository,
) authApp.UserCreatedUseCase {
	return authApp.NewUserCreatedUseCase(eventIDGen, eventRepo, userRepo)
}
func provideSocialUserCreatedUseCase(userRepo *socialPersistence.UserRepository) socialApp.UserCreatedUseCase {
	return socialApp.NewUserCreatedUseCase(userRepo)
}
func provideSocialRoomCreatedUseCase(
	eventIDGen event.IDGenerator,
	roomRepo socialDomain.RoomRepository,
	eventRepo event.Repository,
) socialApp.RoomCreatedUseCase {
	return socialApp.NewRoomCreatedUseCase(eventIDGen, eventRepo, roomRepo)
}
func provideSocialRoomJoinedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) socialApp.RoomJoinedUseCase {
	return socialApp.NewRoomJoinedUseCase(eventIDGen, eventRepo)
}
func provideSocialRoomLeftUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) socialApp.RoomLeftUseCase {
	return socialApp.NewRoomLeftUseCase(eventIDGen, eventRepo)
}
func provideChatUserCreatedUseCase(userRepo chatDomain.UserRepository) chatApp.UserCreatedUseCase {
	return chatApp.NewUserCreatedUseCase(userRepo)
}
func provideChatRoomCreatedUseCase(
	roomRepo chatDomain.RoomRepository,
) chatApp.RoomCreatedUseCase {
	return chatApp.NewRoomCreatedUseCase(roomRepo)
}
func provideChatRoomJoinedUseCase(roomRepo chatDomain.RoomRepository) chatApp.RoomJoinedUseCase {
	return chatApp.NewRoomJoinedUseCase(roomRepo, roomRepo)
}
func provideChatRoomLeftUseCase(roomRepo chatDomain.RoomRepository) chatApp.RoomLeftUseCase {
	return chatApp.NewRoomLeftUseCase(roomRepo, roomRepo)
}
func provideChatPrivateMessageCreatedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
	messageRepo chatDomain.PrivateMessageRepository,
) chatApp.PrivateMessageCreatedUseCase {
	return chatApp.NewPrivateMessageCreatedUseCase(eventIDGen, messageRepo, eventRepo)
}
func provideChatRoomMessageCreatedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
	roomRepo chatDomain.RoomRepository,
	messageRepo chatDomain.RoomMessageRepository,
) chatApp.RoomMessageCreatedUseCase {
	return chatApp.NewRoomMessageCreatedUseCase(eventIDGen, messageRepo, roomRepo, eventRepo)
}
func provideNotificationWelcomeEmailNotificationRequestedUseCase(emailNotifier notificationApp.WelcomeEmailNotifier) notificationApp.WelcomeEmailNotificationRequestedUseCase {
	return notificationApp.NewWelcomeEmailNotificationRequestedUseCase(emailNotifier)
}
func provideNotificationPrivateMessageNotificationRequestedUseCase(
	messageRepo notificationDomain.PrivateMessageRepository,
	messageNotifier notificationDomain.PrivateMessageNotifier,
) notificationApp.PrivateMessageNotificationRequestedUseCase {
	return notificationApp.NewPrivateMessageNotificationRequestedUseCase(messageRepo, messageNotifier)
}
func provideNotificationRoomMessageNotificationRequestedUseCase(
	messageRepo notificationDomain.RoomMessageRepository,
	messageNotifier notificationDomain.RoomMessageNotifier,
) notificationApp.RoomMessageNotificationRequestedUseCase {
	return notificationApp.NewRoomMessageNotificationRequestedUseCase(messageNotifier, messageRepo)
}
func provideNotificationUndeliveredMessagesNotificationRequestedUseCase(
	privateMessageRepo notificationDomain.PrivateMessageRepository,
	roomMessageRepo notificationDomain.RoomMessageRepository,
	privateMessageNotifier notificationDomain.PrivateMessageNotifier,
	roomMessageNotifier notificationDomain.RoomMessageNotifier,
) notificationApp.UndeliveredMessagesNotificationRequestedUseCase {
	return notificationApp.NewUndeliveredMessagesNotificationRequestedUseCase(privateMessageRepo, privateMessageRepo, privateMessageNotifier, roomMessageRepo, roomMessageRepo, roomMessageNotifier)
}
