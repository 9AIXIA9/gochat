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
	"gochat/internal/roomship/domain"
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
	provideRefuseMemberRequestUseCase,
	provideAgreeMemberRequestUseCase,
	provideSendMemberRequestUseCase,
	provideCreateRoomUseCase,
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
	provideRoomshipRoomshipCreatedUseCase,
	provideRoomshipMemberRequestAgreedUseCase,
	provideRoomshipMemberRequestCreatedUseCase,
	provideRoomshipRoomCreatedUseCase,
	provideChatUserCreatedUseCase,
	provideChatRoomCreatedUseCase,
	provideChatRoomMessageCreatedUseCase,
	provideChatPrivateMessageCreatedUseCase,
	provideNotificationWelcomeEmailNotificationRequestedUseCase,
	provideNotificationRoomMessageNotificationRequestedUseCase,
	provideNotificationSystemMessageNotificationRequestedUseCase,
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
	messageRepo chatDomain.PrivateMessageRepository,
	friendshipRepo chatDomain.FriendshipRepository,
) chatApp.SendPrivateMessageUseCase {
	return chatApp.NewSendPrivateMessageUseCase(friendshipRepo, messageIDGen, eventIDGen, messageRepo)
}
func provideSendRoomMessageUseCase(
	messageIDGen kernel.MessageIDGenerator,
	eventIDGen event.IDGenerator,
	roomshipRepo chatDomain.RoomshipRepository,
	messageRepo chatDomain.RoomMessageRepository,
) chatApp.SendRoomMessageUseCase {
	return chatApp.NewSendRoomMessageUseCase(messageIDGen, eventIDGen, roomshipRepo, messageRepo)
}
func provideCreateRoomUseCase(
	encryptor roomshipDomain.Encryptor,
	roomRepo roomshipDomain.RoomRepository,
	eventIDGenerator event.IDGenerator,
	roomIDGenerator roomshipDomain.RoomIDGenerator,
	roomNumberGenerator roomshipDomain.RoomNumberGenerator,
) (roomshipApp.CreateRoomUseCase, error) {
	return roomshipApp.NewCreateRoomUseCase(
		eventIDGenerator,
		roomIDGenerator,
		roomNumberGenerator,
		encryptor,
		roomRepo,
	)
}
func provideSendMemberRequestUseCase(
	roomshipRepo roomshipDomain.RoomshipRepository,
	requestRepo roomshipDomain.MemberRequestRepository,
	roomRepo roomshipDomain.RoomRepository,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator kernel.OperationIDGenerator,
	comparator roomshipDomain.Comparator,
) (roomshipApp.SendMemberRequestUseCase, error) {
	return roomshipApp.NewSendMemberRequestUseCase(
		requestRepo,
		roomshipRepo,
		roomRepo,
		eventIDGenerator,
		operationIDGenerator,
		comparator,
		requestRepo,
	)
}
func provideAgreeMemberRequestUseCase(
	requestRepo roomshipDomain.MemberRequestRepository,
	roomshipRepo roomshipDomain.RoomshipRepository,
	idGenerator event.IDGenerator,
) (roomshipApp.AgreeMemberRequestUseCase, error) {
	return roomshipApp.NewAgreeMemberRequestUseCase(
		requestRepo,
		roomshipRepo,
		requestRepo,
		idGenerator,
	)
}
func provideRefuseMemberRequestUseCase(
	requestRepo roomshipDomain.MemberRequestRepository,
	roomshipRepo roomshipDomain.RoomshipRepository,
) (roomshipApp.RefuseMemberRequestUseCase, error) {
	return roomshipApp.NewRefuseMemberRequestUseCase(
		requestRepo,
		roomshipRepo,
		requestRepo,
	)
}
func provideSendFriendRequestUseCase(
	friendshipRepo friendshipDomain.FriendshipRepository,
	requestRepo friendshipDomain.FriendRequestRepository,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator kernel.OperationIDGenerator,
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
func provideRoomshipUserCreatedUseCase(userRepo *RoomshipPersistence.UserRepository) (roomshipApp.UserCreatedUseCase, error) {
	return roomshipApp.NewUserCreatedUseCase(userRepo)
}
func provideRoomshipRoomCreatedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	roomRepo domain.RoomRepository,
	roomshipRepo domain.RoomshipRepository,
) (roomshipApp.RoomCreatedUseCase, error) {
	return roomshipApp.NewRoomCreatedUseCase(
		roomshipIDGenerator,
		idGenerator,
		roomRepo,
		roomshipRepo,
	)
}
func provideRoomshipMemberRequestCreatedUseCase(
	requestRepo domain.MemberRequestRepository,
	roomshipRepo domain.RoomshipRepository,
	idGenerator event.IDGenerator,
	eventRepo event.Repository,
) (roomshipApp.MemberRequestCreatedUseCase, error) {
	return roomshipApp.NewMemberRequestCreatedUseCase(
		requestRepo, roomshipRepo, idGenerator, eventRepo,
	)
}
func provideRoomshipMemberRequestAgreedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	requestRepo domain.MemberRequestRepository,
	roomshipRepo domain.RoomshipRepository,
) (roomshipApp.MemberRequestAgreedUseCase, error) {
	return roomshipApp.NewMemberRequestAgreedUseCase(
		roomshipIDGenerator,
		idGenerator,
		requestRepo,
		roomshipRepo,
	)
}
func provideRoomshipRoomshipCreatedUseCase(
	roomshipRepo domain.RoomshipRepository,
	eventRepo event.Repository,
	idGenerator event.IDGenerator,
) (roomshipApp.RoomshipCreatedUseCase, error) {
	return roomshipApp.NewRoomshipCreatedUseCase(
		roomshipRepo,
		roomshipRepo,
		idGenerator,
		eventRepo,
	)
}
func provideChatUserCreatedUseCase(userRepo chatDomain.UserRepository) chatApp.UserCreatedUseCase {
	return chatApp.NewUserCreatedUseCase(userRepo)
}
func provideChatRoomCreatedUseCase(
	roomRepo chatDomain.RoomRepository,
) chatApp.RoomCreatedUseCase {
	return chatApp.NewRoomCreatedUseCase(roomRepo)
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
	roomshipRepo chatDomain.RoomshipRepository,
	messageRepo chatDomain.RoomMessageRepository,
) chatApp.RoomMessageCreatedUseCase {
	return chatApp.NewRoomMessageCreatedUseCase(eventIDGen, messageRepo, roomshipRepo, eventRepo)
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
func provideNotificationSystemMessageNotificationRequestedUseCase(
	idGenerator kernel.MessageIDGenerator,
	messageRepo notificationDomain.SystemMessageRepository,
	messageNotifier notificationDomain.SystemMessageNotifier,
) (notificationApp.SystemMessageNotificationRequestedUseCase, error) {
	return notificationApp.NewSystemMessageNotificationRequestedUseCase(
		messageRepo,
		messageNotifier,
		idGenerator,
	)
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
