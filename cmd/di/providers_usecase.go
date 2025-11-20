package di

import (
	"gochat/internal/application/usecase"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	chatApp "gochat/internal/chat/application"
	chatUsecase "gochat/internal/chat/application/usecase"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/uuid"
	notificationUsecase "gochat/internal/notification/application/usecase"
	gomailUtil "gochat/internal/notification/infrastructure/gomail"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	notificationWebsocketInfrastructure "gochat/internal/notification/infrastructure/websocket"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	application2 "gochat/internal/social/application"
	domain2 "gochat/internal/social/domain"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"

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
	provideNotificationUndeliveredMessageNotificationRequestedUseCase,
	provideNotificationMessageReadUseCase,
)

var UseCaseKafkaSet = wire.NewSet(
	provideWebsocketUserSessionStartedUseCase,
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
	provideNotificationMessageNotificationRequestedUseCase,
)

// -------------------- UseCases (HTTP side) --------------------
func provideSignUpUseCase(
	eventIDGen event.IDGenerator,
	userIDGen domain.UserIDGenerator,
	numberGen domain.UserNumberGenerator,
	encryptor domain.Encryptor,
	userSaver domain.UserSaver,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) application.SignUpUseCase {
	return application.NewSignUpUseCase(eventIDGen, userIDGen, numberGen, encryptor, userSaver, eventSaver, unitOfWork)
}

func provideLoginUseCase(
	eventIDGen event.IDGenerator,
	comparator domain.Comparator,
	userFinder domain.UserFinderByNumber,
	refreshTokenSaver domain.RefreshTokenSaver,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) application.LoginUseCase {
	return application.NewLoginUseCase(eventIDGen, comparator, userFinder, refreshTokenSaver, accessTokenGenerator, refreshTokenGenerator)
}
func provideRefreshAccessTokenUseCase(
	refreshTokenSaver domain.RefreshTokenSaver,
	refreshTokenFinder domain.RefreshTokenFinder,
	accessTokenGenerator domain.AccessTokenGenerator,
	refreshTokenGenerator domain.RefreshTokenGenerator,
) application.RefreshAccessTokenUseCase {
	return application.NewRefreshAccessTokenUseCase(refreshTokenSaver, refreshTokenFinder, accessTokenGenerator, refreshTokenGenerator)
}
func provideParseAccessTokenUseCase(accessTokenParser domain.AccessTokenParser) application.ParseAccessTokenUseCase {
	return application.NewParseAccessTokenUseCase(accessTokenParser)
}
func provideSendPrivateMessageUseCase(
	messageIDGen chatApp.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	userRepo *chatRepository.UserRepository,
	msgRepo *chatRepository.MessageRepository,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) chatUsecase.SendPrivateMessageUseCase {
	return chatUsecase.NewSendPrivateMessageUseCase(messageIDGen, eventIDGen, userRepo, msgRepo, eventSaver, unitOfWork)
}
func provideSendRoomMessageUseCase(
	messageIDGen chatApp.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomRepo *chatRepository.RoomRepository,
	msgRepo *chatRepository.MessageRepository,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) chatUsecase.SendRoomMessageUseCase {
	return chatUsecase.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomRepo, msgRepo, eventSaver, unitOfWork)
}
func provideCreateRoomUseCase(
	eventIDGen event.IDGenerator,
	roomIDGen domain2.RoomIDGenerator,
	numberGen domain2.RoomNumberGenerator,
	encryptor domain2.Encryptor,
	roomSaver domain2.RoomSaver,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) application2.CreateRoomUseCase {
	return application2.NewCreateRoomUseCase(eventIDGen, roomIDGen, numberGen, encryptor, roomSaver, eventSaver, unitOfWork)
}
func provideJoinRoomUseCase(
	eventIDGen event.IDGenerator,
	eventSaver event.UnpublishedEventsSaver,
	finder domain2.RoomFinderByNumber,
	comparator domain2.Comparator,
	roomMemberSaver domain2.RoomMemberSaver,
	unitOfWork kernel.UnitOfWork,
) application2.JoinRoomUseCase {
	return application2.NewJoinRoomUseCase(eventIDGen, eventSaver, finder, comparator, roomMemberSaver, unitOfWork)
}
func provideLeaveRoomUseCase(
	eventIDGen event.IDGenerator,
	finder domain2.RoomFinderByNumber,
	roomMemberDeleter domain2.RoomMemberDeleter,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) application2.LeaveRoomUseCase {
	return application2.NewLeaveRoomUseCase(eventIDGen, finder, roomMemberDeleter, eventSaver, unitOfWork)
}

// -------------------- Event UseCases (Kafka consumer side) --------------------
func provideWebsocketUserSessionStartedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver) usecase.UserSessionStartedUseCase {
	return usecase.NewUserSessionStartedUseCase(eventIDGen, saver)
}
func provideAuthUserCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventsSaver, userRepo *authorizationRepository.UserRepository) application.UserCreatedUseCase {
	return application.NewUserCreatedUseCase(eventIDGen, saver, userRepo)
}
func provideSocialUserCreatedUseCase(userRepo *socialRepository.UserRepository) application2.UserCreatedUseCase {
	return application2.NewUserCreatedUseCase(userRepo)
}
func provideSocialRoomCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver, roomRepo *socialRepository.RoomRepository) application2.RoomCreatedUseCase {
	return application2.NewRoomCreatedUseCase(eventIDGen, saver, roomRepo)
}
func provideSocialRoomJoinedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver) application2.RoomJoinedUseCase {
	return application2.NewRoomJoinedUseCase(eventIDGen, saver)
}
func provideSocialRoomLeftUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver) application2.RoomLeftUseCase {
	return application2.NewRoomLeftUseCase(eventIDGen, saver)
}
func provideChatUserCreatedUseCase(userRepo *chatRepository.UserRepository) chatUsecase.UserCreatedUseCase {
	return chatUsecase.NewUserCreatedUseCase(userRepo)
}
func provideChatRoomCreatedUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomCreatedUseCase {
	return chatUsecase.NewRoomCreatedUseCase(roomRepo)
}
func provideChatRoomJoinedUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomJoinedUseCase {
	return chatUsecase.NewRoomJoinedUseCase(roomRepo)
}
func provideChatRoomLeftUseCase(roomRepo *chatRepository.RoomRepository) chatUsecase.RoomLeftUseCase {
	return chatUsecase.NewRoomLeftUseCase(roomRepo)
}
func provideChatPrivateMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver, messageRepo *chatRepository.MessageRepository) chatUsecase.PrivateMessageCreatedUseCase {
	return chatUsecase.NewPrivateMessageCreatedUseCase(eventIDGen, messageRepo, saver)
}
func provideChatRoomMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventsSaver, messageRepo *chatRepository.MessageRepository) chatUsecase.RoomMessageCreatedUseCase {
	return chatUsecase.NewRoomMessageCreatedUseCase(eventIDGen, messageRepo, saver)
}
func provideNotificationWelcomeEmailNotificationRequestedUseCase(emailNotifier *gomailUtil.EmailNotifier) notificationUsecase.WelcomeEmailNotificationRequestedUseCase {
	return notificationUsecase.NewWelcomeEmailNotificationRequestedUseCase(emailNotifier)
}
func provideNotificationMessageNotificationRequestedUseCase(messageRepo *notificationRepository.MessageRepository, messageNotifier *notificationWebsocketInfrastructure.MessageNotifier) notificationUsecase.MessageNotificationRequestedUseCase {
	return notificationUsecase.NewMessageNotificationRequestedUseCase(messageRepo, messageNotifier, messageRepo)
}
func provideNotificationUndeliveredMessageNotificationRequestedUseCase(messageRepo *notificationRepository.MessageRepository, messageNotifier *notificationWebsocketInfrastructure.MessageNotifier) notificationUsecase.UndeliveredMessageNotificationRequestedUseCase {
	return notificationUsecase.NewUndeliveredMessageNotificationRequestedUseCase(messageRepo, messageNotifier, messageRepo)
}
func provideNotificationMessageReadUseCase(messageRepo *notificationRepository.MessageRepository) notificationUsecase.MessageReadUseCase {
	return notificationUsecase.NewMessageReadUseCase(messageRepo)
}
