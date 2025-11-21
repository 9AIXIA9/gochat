package di

import (
	application5 "gochat/internal/application"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	application3 "gochat/internal/chat/application"
	chatApp "gochat/internal/chat/domain"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/uuid"
	application4 "gochat/internal/notification/application"
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
	provideNotificationPrivateMessageReadUseCase,
	provideNotificationRoomMessageReadUseCase,
	provideNotificationReceivedPrivateMessageNotificationRequestedUseCase,
	provideNotificationReceivedRoomMessageNotificationRequestedUseCase,
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
	provideNotificationRoomMessageNotificationRequestedUseCase,
	provideNotificationPrivateMessageNotificationRequestedUseCase,
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
	msgRepo *chatRepository.PrivateMessageRepository,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) application3.SendPrivateMessageUseCase {
	return application3.NewSendPrivateMessageUseCase(messageIDGen, eventIDGen, userRepo, msgRepo, eventSaver, unitOfWork)
}
func provideSendRoomMessageUseCase(
	messageIDGen chatApp.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomRepo *chatRepository.RoomRepository,
	msgRepo *chatRepository.RoomMessageRepository,
	eventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) application3.SendRoomMessageUseCase {
	return application3.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomRepo, msgRepo, eventSaver, unitOfWork)
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
func provideWebsocketUserSessionStartedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventsSaver) application5.UserSessionStartedUseCase {
	return application5.NewUserSessionStartedUseCase(eventIDGen, saver)
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
func provideChatUserCreatedUseCase(userRepo *chatRepository.UserRepository) application3.UserCreatedUseCase {
	return application3.NewUserCreatedUseCase(userRepo)
}
func provideChatRoomCreatedUseCase(roomRepo *chatRepository.RoomRepository, unitOfWork *gormutils.UnitOfWork) application3.RoomCreatedUseCase {
	return application3.NewRoomCreatedUseCase(roomRepo, roomRepo, unitOfWork)
}
func provideChatRoomJoinedUseCase(roomRepo *chatRepository.RoomRepository) application3.RoomJoinedUseCase {
	return application3.NewRoomJoinedUseCase(roomRepo)
}
func provideChatRoomLeftUseCase(roomRepo *chatRepository.RoomRepository) application3.RoomLeftUseCase {
	return application3.NewRoomLeftUseCase(roomRepo)
}
func provideChatPrivateMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver, messageRepo *chatRepository.PrivateMessageRepository) application3.PrivateMessageCreatedUseCase {
	return application3.NewPrivateMessageCreatedUseCase(eventIDGen, messageRepo, saver)
}
func provideChatRoomMessageCreatedUseCase(eventIDGen *uuid.EventIDGenerator, saver event.UnpublishedEventSaver, roomRepo *chatRepository.RoomRepository, messageRepo *chatRepository.RoomMessageRepository) application3.RoomMessageCreatedUseCase {
	return application3.NewRoomMessageCreatedUseCase(eventIDGen, messageRepo, roomRepo, saver)
}
func provideNotificationWelcomeEmailNotificationRequestedUseCase(emailNotifier *gomailUtil.EmailNotifier) application4.WelcomeEmailNotificationRequestedUseCase {
	return application4.NewWelcomeEmailNotificationRequestedUseCase(emailNotifier)
}
func provideNotificationPrivateMessageNotificationRequestedUseCase(messageRepo *notificationRepository.PrivateMessageRepository, messageNotifier *notificationWebsocketInfrastructure.PrivateMessageNotifier) application4.PrivateMessageNotificationRequestedUseCase {
	return application4.NewPrivateMessageNotificationRequestedUseCase(messageRepo, messageNotifier)
}
func provideNotificationRoomMessageNotificationRequestedUseCase(messageRepo *notificationRepository.RoomMessageRepository, messageNotifier *notificationWebsocketInfrastructure.RoomMessageNotifier) application4.RoomMessageNotificationRequestedUseCase {
	return application4.NewRoomMessageNotificationRequestedUseCase(messageNotifier, messageRepo)
}
func provideNotificationReceivedPrivateMessageNotificationRequestedUseCase(messageRepo *notificationRepository.PrivateMessageRepository, messageNotifier *notificationWebsocketInfrastructure.PrivateMessageNotifier) application4.ReceivedPrivateMessageNotificationRequestedUseCase {
	return application4.NewReceivedPrivateMessageNotificationRequestedUseCase(messageRepo, messageRepo, messageNotifier)
}
func provideNotificationReceivedRoomMessageNotificationRequestedUseCase(messageRepo *notificationRepository.RoomMessageRepository, messageNotifier *notificationWebsocketInfrastructure.RoomMessageNotifier) application4.ReceivedRoomMessageNotificationRequestedUseCase {
	return application4.NewReceivedRoomMessageNotificationRequestedUseCase(messageRepo, messageRepo, messageNotifier)
}
func provideNotificationPrivateMessageReadUseCase(messageRepo *notificationRepository.PrivateMessageRepository) application4.PrivateMessageReadUseCase {
	return application4.NewPrivateMessageReadUseCase(messageRepo, messageRepo)
}
func provideNotificationRoomMessageReadUseCase(messageRepo *notificationRepository.RoomMessageRepository) application4.RoomMessageReadUseCase {
	return application4.NewRoomMessageReadUseCase(messageRepo, messageRepo)
}
