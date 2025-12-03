package di

import (
	rootapp "gochat/internal/application"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	friendshipApp "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	roomshipApp "gochat/internal/roomship/application"
	roomshipDomain "gochat/internal/roomship/domain"
	RoomshipPersistence "gochat/internal/roomship/infrastructure/persistence/repository"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

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
	provideListFriendRequestsUseCase,
	provideRefuseFriendRequestUseCase,
	provideAgreeFriendRequestUseCase,
	provideSendFriendRequestUseCase,
)

var UseCaseWebsocketSet = wire.NewSet(
	provideWebsocketUserSessionStartedUseCase,
	provideNotificationReadPrivateMessageUseCase,
	provideNotificationReadRoomMessageUseCase,
)

var UseCaseKafkaSet = wire.NewSet(
	provideAuthUserCreatedUseCase,
	provideRoomshipUserCreatedUseCase,
	provideRoomshipRoomCreatedUseCase,
	provideRoomshipRoomJoinedUseCase,
	provideRoomshipRoomLeftUseCase,
	provideChatUserCreatedUseCase,
	provideChatRoomCreatedUseCase,
	provideChatRoomJoinedUseCase,
	provideChatRoomLeftUseCase,
	provideChatRoomMessageCreatedUseCase,
	provideChatPrivateMessageCreatedUseCase,
	provideNotificationWelcomeEmailNotificationRequestedUseCase,
	provideNotificationRoomMessageNotificationRequestedUseCase,
	provideNotificationFriendRequestCreatedNotificationRequestedUseCase,
	provideNotificationFriendshipCreatedNotificationRequestedUseCase,
	provideNotificationPrivateMessageNotificationRequestedUseCase,
	provideNotificationUndeliveredMessagesNotificationRequestedUseCase,
	provideFriendshipUserCreatedUseCase,
	provideFriendshipFriendshipCreatedUseCase,
	provideFriendshipFriendRequestCreatedUseCase,
	provideFriendshipFriendRequestAgreedUseCase,
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
	messageIDGen kernel.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	userRepo chatDomain.UserRepository,
	messageRepo chatDomain.PrivateMessageRepository,
) chatApp.SendPrivateMessageUseCase {
	return chatApp.NewSendPrivateMessageUseCase(messageIDGen, eventIDGen, userRepo, messageRepo)
}
func provideSendRoomMessageUseCase(
	messageIDGen kernel.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomRepo chatDomain.RoomRepository,
	messageRepo chatDomain.RoomMessageRepository,
) chatApp.SendRoomMessageUseCase {
	return chatApp.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomRepo, messageRepo)
}
func provideCreateRoomUseCase(
	eventIDGen event.IDGenerator,
	roomIDGen roomshipDomain.RoomIDGenerator,
	numberGen roomshipDomain.RoomNumberGenerator,
	encryptor roomshipDomain.Encryptor,
	roomRepo roomshipDomain.RoomRepository,
) roomshipApp.CreateRoomUseCase {
	return roomshipApp.NewCreateRoomUseCase(eventIDGen, roomIDGen, numberGen, encryptor, roomRepo)
}
func provideJoinRoomUseCase(
	eventIDGen event.IDGenerator,
	comparator roomshipDomain.Comparator,
	roomRepo roomshipDomain.RoomRepository,
) roomshipApp.JoinRoomUseCase {
	return roomshipApp.NewJoinRoomUseCase(eventIDGen, roomRepo, comparator, roomRepo)
}
func provideLeaveRoomUseCase(
	eventIDGen event.IDGenerator,
	roomRepo roomshipDomain.RoomRepository,
) roomshipApp.LeaveRoomUseCase {
	return roomshipApp.NewLeaveRoomUseCase(eventIDGen, roomRepo, roomRepo)
}
func provideSendFriendRequestUseCase(
	friendshipRepo friendshipDomain.FriendshipRepository,
	requestRepo friendshipDomain.FriendRequestRepository,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator friendshipDomain.OperationIDGenerator,
) (friendshipApp.SendFriendRequestUseCase, error) {
	return friendshipApp.NewSendFriendRequestUseCase(
		friendshipRepo,
		requestRepo,
		requestRepo,
		eventIDGenerator,
		operationIDGenerator,
	)
}
func provideAgreeFriendRequestUseCase(
	requestRepo friendshipDomain.FriendRequestRepository,
	idGenerator event.IDGenerator,
) (friendshipApp.AgreeFriendRequestUseCase, error) {
	return friendshipApp.NewAgreeFriendRequestUseCase(
		requestRepo,
		requestRepo,
		idGenerator,
	)
}
func provideRefuseFriendRequestUseCase(
	requestRepo friendshipDomain.FriendRequestRepository,
) (friendshipApp.RefuseFriendRequestUseCase, error) {
	return friendshipApp.NewRefuseFriendRequestUseCase(
		requestRepo,
		requestRepo,
	)
}
func provideListFriendRequestsUseCase(
	friendRequestRepo friendshipDomain.FriendRequestRepository,
) (friendshipApp.ListFriendRequestsUseCase, error) {
	return friendshipApp.NewListFriendRequestsUseCase(
		friendRequestRepo,
	)
}

// -------------------- Event UseCases (websocket side) --------------------
func provideWebsocketUserSessionStartedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) rootapp.UserSessionStartedUseCase {
	return rootapp.NewUserSessionStartedUseCase(eventIDGen, eventRepo)
}
func provideNotificationReadPrivateMessageUseCase(messageRepo notificationDomain.PrivateMessageRepository) notificationApp.ReadPrivateMessageUseCase {
	return notificationApp.NewReadPrivateMessageUseCase(messageRepo, messageRepo)
}
func provideNotificationReadRoomMessageUseCase(messageRepo notificationDomain.RoomMessageRepository) notificationApp.ReadRoomMessageUseCase {
	return notificationApp.NewReadRoomMessageUseCase(messageRepo, messageRepo)
}

// -------------------- Event UseCases (Kafka consumer side) --------------------

func provideAuthUserCreatedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
	userRepo authDomain.UserRepository,
) authApp.UserCreatedUseCase {
	return authApp.NewUserCreatedUseCase(eventIDGen, eventRepo, userRepo)
}
func provideRoomshipUserCreatedUseCase(userRepo *RoomshipPersistence.UserRepository) roomshipApp.UserCreatedUseCase {
	return roomshipApp.NewUserCreatedUseCase(userRepo)
}
func provideRoomshipRoomCreatedUseCase(
	eventIDGen event.IDGenerator,
	roomRepo roomshipDomain.RoomRepository,
	eventRepo event.Repository,
) roomshipApp.RoomCreatedUseCase {
	return roomshipApp.NewRoomCreatedUseCase(eventIDGen, eventRepo, roomRepo)
}
func provideRoomshipRoomJoinedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) roomshipApp.RoomJoinedUseCase {
	return roomshipApp.NewRoomJoinedUseCase(eventIDGen, eventRepo)
}
func provideRoomshipRoomLeftUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) roomshipApp.RoomLeftUseCase {
	return roomshipApp.NewRoomLeftUseCase(eventIDGen, eventRepo)
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
func provideNotificationFriendRequestCreatedNotificationRequestedUseCase(
	idGenerator kernel.MessageIDGenerator,
	messageRepo notificationDomain.SystemMessageRepository,
	messageNotifier notificationDomain.SystemMessageNotifier,
) (notificationApp.FriendRequestCreatedNotificationRequestedUseCase, error) {
	return notificationApp.NewFriendRequestCreatedNotificationRequestedUseCase(idGenerator, messageRepo, messageNotifier)
}
func provideNotificationFriendshipCreatedNotificationRequestedUseCase(
	idGenerator kernel.MessageIDGenerator,
	messageRepo notificationDomain.SystemMessageRepository,
	messageNotifier notificationDomain.SystemMessageNotifier,
) (notificationApp.FriendshipCreatedNotificationRequestedUseCase, error) {
	return notificationApp.NewFriendshipCreatedNotificationRequestedUseCase(idGenerator, messageRepo, messageNotifier)
}
func provideNotificationRoomMessageNotificationRequestedUseCase(
	messageRepo notificationDomain.RoomMessageRepository,
	messageNotifier notificationDomain.RoomMessageNotifier,
) notificationApp.RoomMessageNotificationRequestedUseCase {
	return notificationApp.NewRoomMessageNotificationRequestedUseCase(messageNotifier, messageRepo)
}
func provideNotificationUndeliveredMessagesNotificationRequestedUseCase(
	systemMessageRepo notificationDomain.SystemMessageRepository,
	systemMessageNotifier notificationDomain.SystemMessageNotifier,
	privateMessageRepo notificationDomain.PrivateMessageRepository,
	privateMessageNotifier notificationDomain.PrivateMessageNotifier,
	roomMessageRepo notificationDomain.RoomMessageRepository,
	roomMessageNotifier notificationDomain.RoomMessageNotifier,
) notificationApp.UndeliveredMessagesNotificationRequestedUseCase {
	return notificationApp.NewUndeliveredMessagesNotificationRequestedUseCase(
		systemMessageRepo,
		systemMessageRepo,
		systemMessageNotifier,
		privateMessageRepo,
		privateMessageRepo,
		privateMessageNotifier,
		roomMessageRepo,
		roomMessageRepo,
		roomMessageNotifier,
	)
}
func provideFriendshipUserCreatedUseCase(
	userRepo friendshipDomain.UserRepository,
) (friendshipApp.UserCreatedUseCase, error) {
	return friendshipApp.NewUserCreatedUseCase(userRepo)
}
func provideFriendshipFriendRequestAgreedUseCase(
	requestRepo friendshipDomain.FriendRequestRepository,
	friendshipRepo friendshipDomain.FriendshipRepository,
	friendshipIDGen friendshipDomain.FriendshipIDGenerator,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendRequestAgreedUseCase, error) {
	return friendshipApp.NewFriendRequestAgreedUseCase(requestRepo, friendshipRepo, friendshipIDGen, eventIDGen)
}
func provideFriendshipFriendRequestCreatedUseCase(
	friendRequestRepo friendshipDomain.FriendRequestRepository,
	eventRepo event.Repository,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendRequestCreatedUseCase, error) {
	return friendshipApp.NewFriendRequestCreatedUseCase(friendRequestRepo, eventRepo, eventIDGen)
}
func provideFriendshipFriendshipCreatedUseCase(
	friendshipRepo friendshipDomain.FriendshipRepository,
	eventRepo event.Repository,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendshipCreatedUseCase, error) {
	return friendshipApp.NewFriendshipCreatedUseCase(friendshipRepo, eventRepo, eventIDGen)
}
