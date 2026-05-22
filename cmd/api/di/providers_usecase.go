package di

import (
	"gochat/config"
	rootapp "gochat/internal/application"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	friendshipApp "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	profileApp "gochat/internal/profile/application"
	profileDomain "gochat/internal/profile/domain"
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
	provideProfileGetRoomProfileUseCase,
	provideProfileGetUserProfileUseCase,
	provideProfileUpdateRoomProfileUseCase,
	provideProfileUpdateUserProfileUseCase,
	provideSendPrivateMessageUseCase,
	provideSendRoomMessageUseCase,
	provideRefuseMemberRequestUseCase,
	provideAgreeMemberRequestUseCase,
	provideSendMemberRequestUseCase,
	provideLeaveRoomUseCase,
	provideListRoomMembersUseCase,
	provideListRoomMessagesUseCase,
	provideListPrivateMessagesUseCase,
	provideCreateRoomUseCase,
	provideListFriendRequestsUseCase,
	provideListFriendshipsUseCase,
	provideRefuseFriendRequestUseCase,
	provideAgreeFriendRequestUseCase,
	provideSendFriendRequestUseCase,
	provideListRoomshipsUseCase,
	provideListMemberRequestsUseCase,
	provideListSystemMessagesUseCase,
)

var UseCaseWebsocketSet = wire.NewSet(
	provideWebsocketUserSessionStartedUseCase,
	provideChatReadRoomMessagesUseCase,
	provideChatReadPrivateMessagesUseCase,
	provideChatConfirmPrivateMessagesUseCase,
	provideChatConfirmRoomMessagesUseCase,
	provideNotificationConfirmSystemMessagesUseCase,
)

var UseCaseKafkaSet = wire.NewSet(
	provideAuthUserCreatedUseCase,
	provideUnpublishedEventsCreatedCase,
	provideProfileRoomshipCreatedUseCase,
	provideProfileRoomCreatedUseCase,
	provideProfileUserCreatedUseCase,
	provideRoomshipUserCreatedUseCase,
	provideRoomshipRoomshipCreatedUseCase,
	provideRoomshipMemberRequestAgreedUseCase,
	provideRoomshipMemberRequestCreatedUseCase,
	provideRoomshipRoomCreatedUseCase,
	provideChatUserCreatedUseCase,
	provideChatRoomCreatedUseCase,
	provideChatRoomshipCreatedUseCase,
	provideChatFriendshipCreatedUseCase,
	provideChatUndeliveredMessagesPushRequestedUseCase,
	provideNotificationWelcomeEmailNotificationRequestedUseCase,
	provideNotificationSystemMessageNotificationRequestedUseCase,
	provideNotificationUndeliveredMessagesNotificationRequestedUseCase,
	provideFriendshipUserCreatedUseCase,
	provideFriendshipFriendshipCreatedUseCase,
	provideFriendshipFriendRequestCreatedUseCase,
	provideFriendshipFriendRequestAgreedUseCase,
)

// -------------------- UseCases (HTTP side) --------------------
func provideSignUpUseCase(
	eventIDGen event.IDGenerator,
	userIDGen authDomain.UserIDGenerator,
	numberGen authDomain.UserNumberGenerator,
	encryptor authDomain.Encryptor,
	userRepo authDomain.UserRepository,
) (authApp.SignUpUseCase, error) {
	return authApp.NewSignUpUseCase(
		eventIDGen,
		userIDGen,
		numberGen,
		encryptor,
		userRepo,
	)
}

func provideLoginUseCase(
	comparator authDomain.Comparator,
	accessTokenGenerator authDomain.AccessTokenGenerator,
	refreshTokenGenerator authDomain.RefreshTokenGenerator,
	userRepo authDomain.UserRepository,
	refreshTokenRepo authDomain.RefreshTokenRepository,
) (authApp.LoginUseCase, error) {
	return authApp.NewLoginUseCase(
		comparator,
		userRepo,
		refreshTokenRepo,
		accessTokenGenerator,
		refreshTokenGenerator,
	)
}
func provideRefreshAccessTokenUseCase(
	refreshTokenRepo authDomain.RefreshTokenRepository,
	accessTokenGenerator authDomain.AccessTokenGenerator,
	refreshTokenGenerator authDomain.RefreshTokenGenerator,
) (authApp.RefreshAccessTokenUseCase, error) {
	return authApp.NewRefreshAccessTokenUseCase(
		refreshTokenRepo,
		refreshTokenRepo,
		accessTokenGenerator,
		refreshTokenGenerator,
	)
}
func provideParseAccessTokenUseCase(
	accessTokenParser authDomain.AccessTokenParser,
) (authApp.ParseAccessTokenUseCase, error) {
	return authApp.NewParseAccessTokenUseCase(
		accessTokenParser,
	)
}
func provideProfileGetUserProfileUseCase(
	profileRepo profileDomain.UserProfileRepository,
) (profileApp.GetUserProfileUseCase, error) {
	return profileApp.NewGetUserProfileUseCase(
		profileRepo,
	)
}
func provideProfileGetRoomProfileUseCase(
	profileRepo profileDomain.RoomProfileRepository,
) (profileApp.GetRoomProfileUseCase, error) {
	return profileApp.NewGetRoomProfileUseCase(
		profileRepo,
	)
}
func provideProfileUpdateUserProfileUseCase(
	profileRepo profileDomain.UserProfileRepository,
) (profileApp.UpdateUserProfileUseCase, error) {
	return profileApp.NewUpdateUserProfileUseCase(
		profileRepo,
		profileRepo,
	)
}
func provideProfileUpdateRoomProfileUseCase(
	roomshipRepo profileDomain.RoomshipRepository,
	profileRepo profileDomain.RoomProfileRepository,
) (profileApp.UpdateRoomProfileUseCase, error) {
	return profileApp.NewUpdateRoomProfileUseCase(
		profileRepo,
		profileRepo,
		roomshipRepo,
	)
}
func provideSendPrivateMessageUseCase(
	messageIDGen kernel.MessageIDGenerator,
	notifier chatDomain.PrivateMessageNotifier,
	messageRepo chatDomain.PrivateMessageRepository,
	friendshipRepo chatDomain.FriendshipRepository,
) (chatApp.SendPrivateMessageUseCase, error) {
	return chatApp.NewSendPrivateMessageUseCase(
		friendshipRepo,
		messageIDGen,
		notifier,
		messageRepo,
	)
}
func provideSendRoomMessageUseCase(
	messageIDGen kernel.MessageIDGenerator,
	notifier chatDomain.RoomMessageNotifier,
	roomshipRepo chatDomain.RoomshipRepository,
	messageRepo chatDomain.RoomMessageRepository,
) (chatApp.SendRoomMessageUseCase, error) {
	return chatApp.NewSendRoomMessageUseCase(
		messageIDGen,
		notifier,
		roomshipRepo,
		messageRepo,
	)
}
func provideListPrivateMessagesUseCase(
	messageRepo chatDomain.PrivateMessageRepository,
) (chatApp.ListPrivateMessagesUseCase, error) {
	return chatApp.NewListPrivateMessagesUseCase(
		messageRepo,
	)
}
func provideListRoomMessagesUseCase(
	messageRepo chatDomain.RoomMessageRepository,
) (chatApp.ListRoomMessagesUseCase, error) {
	return chatApp.NewListRoomMessagesUseCase(
		messageRepo,
	)
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
func provideLeaveRoomUseCase(
	roomRepo roomshipDomain.RoomRepository,
	roomshipRepo roomshipDomain.RoomshipRepository,
) (roomshipApp.LeaveRoomUseCase, error) {
	return roomshipApp.NewLeaveRoomUseCase(
		roomshipRepo,
		roomRepo,
		roomshipRepo,
		roomshipRepo,
	)
}
func provideListRoomMembersUseCase(
	roomshipRepo roomshipDomain.RoomshipRepository,
) (roomshipApp.ListRoomMembersUseCase, error) {
	return roomshipApp.NewListRoomMembersUseCase(
		roomshipRepo,
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
func provideListMemberRequestsUseCase(
	memberRequestRepo roomshipDomain.MemberRequestRepository,
) (roomshipApp.ListMemberRequestsUseCase, error) {
	return roomshipApp.NewListMemberRequestsUseCase(
		memberRequestRepo,
	)
}
func provideListRoomshipsUseCase(
	roomshipRepo roomshipDomain.RoomshipRepository,
) (roomshipApp.ListRoomshipsUseCase, error) {
	return roomshipApp.NewListRoomshipsUseCase(
		roomshipRepo,
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
func provideListFriendshipsUseCase(
	friendRequestRepo friendshipDomain.FriendshipRepository,
) (friendshipApp.ListFriendshipsUseCase, error) {
	return friendshipApp.NewListFriendshipsUseCase(
		friendRequestRepo,
	)
}
func provideListSystemMessagesUseCase(
	messageRepo notificationDomain.SystemMessageRepository,
) (notificationApp.ListSystemMessagesUseCase, error) {
	return notificationApp.NewListSystemMessagesUseCase(
		messageRepo,
	)
}

// -------------------- Event UseCases (websocket side) --------------------
func provideWebsocketUserSessionStartedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
) (rootapp.UserSessionStartedUseCase, error) {
	return rootapp.NewUserSessionStartedUseCase(
		eventIDGen,
		eventRepo,
	)
}

func provideChatReadPrivateMessagesUseCase(
	messageRepo chatDomain.PrivateMessageRepository,
) (chatApp.ReadPrivateMessagesUseCase, error) {
	return chatApp.NewReadPrivateMessagesUseCase(
		messageRepo,
	)
}

func provideChatReadRoomMessagesUseCase(
	messageRepo chatDomain.RoomMessageRepository,
) (chatApp.ReadRoomMessagesUseCase, error) {
	return chatApp.NewReadRoomMessagesUseCase(
		messageRepo,
	)
}

func provideNotificationConfirmSystemMessagesUseCase(
	updater notificationDomain.SystemMessagesStatesUpdaterByMessageIDs,
) (notificationApp.ConfirmSystemMessagesUseCase, error) {
	return notificationApp.NewConfirmSystemMessagesUseCase(updater)
}

func provideUnpublishedEventsCreatedCase(
	appConfig *config.App,
	eventPublisher event.Publisher,
	eventRepository event.Repository,
) (rootapp.UnpublishedEventsCreatedUseCase, error) {
	_ = appConfig
	return rootapp.NewUnpublishedEventsCreatedUseCase(
		eventPublisher,
		eventRepository,
	)
}

// -------------------- Event UseCases (Kafka consumer side) --------------------

func provideAuthUserCreatedUseCase(
	eventIDGen event.IDGenerator,
	eventRepo event.Repository,
	userRepo authDomain.UserRepository,
) (authApp.UserCreatedUseCase, error) {
	return authApp.NewUserCreatedUseCase(
		eventIDGen,
		eventRepo,
		userRepo,
	)
}
func provideProfileUserCreatedUseCase(
	userRepo profileDomain.UserRepository,
	profileRepo profileDomain.UserProfileRepository,
) (profileApp.UserCreatedUseCase, error) {
	return profileApp.NewUserCreatedUseCase(
		userRepo,
		profileRepo,
	)
}
func provideProfileRoomCreatedUseCase(
	roomRepo profileDomain.RoomRepository,
	profileRepo profileDomain.RoomProfileRepository,
) (profileApp.RoomCreatedUseCase, error) {
	return profileApp.NewRoomCreatedUseCase(
		roomRepo,
		profileRepo,
	)
}
func provideProfileRoomshipCreatedUseCase(
	roomshipRepo profileDomain.RoomshipRepository,
) (profileApp.RoomshipCreatedUseCase, error) {
	return profileApp.NewRoomshipCreatedUseCase(
		roomshipRepo,
	)
}
func provideRoomshipUserCreatedUseCase(
	userRepo *RoomshipPersistence.UserRepository,
) (roomshipApp.UserCreatedUseCase, error) {
	return roomshipApp.NewUserCreatedUseCase(userRepo)
}
func provideRoomshipRoomCreatedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	roomRepo domain.RoomRepository,
	roomshipRepo domain.RoomshipRepository,
	eventRepo event.Repository,
) (roomshipApp.RoomCreatedUseCase, error) {
	return roomshipApp.NewRoomCreatedUseCase(
		roomshipIDGenerator,
		idGenerator,
		roomRepo,
		roomshipRepo,
		eventRepo,
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
func provideChatUserCreatedUseCase(
	userRepo chatDomain.UserRepository,
) (chatApp.UserCreatedUseCase, error) {
	return chatApp.NewUserCreatedUseCase(
		userRepo,
	)
}
func provideChatRoomCreatedUseCase(
	roomRepo chatDomain.RoomRepository,
) (chatApp.RoomCreatedUseCase, error) {
	return chatApp.NewRoomCreatedUseCase(roomRepo)
}
func provideChatFriendshipCreatedUseCase(
	friendshipRepo chatDomain.FriendshipRepository,
) (chatApp.FriendshipCreatedUseCase, error) {
	return chatApp.NewFriendshipCreatedUseCase(
		friendshipRepo,
	)
}
func provideChatRoomshipCreatedUseCase(
	roomshipRepo chatDomain.RoomshipRepository,
) (chatApp.RoomshipCreatedUseCase, error) {
	return chatApp.NewRoomshipCreatedUseCase(
		roomshipRepo,
	)
}
func provideChatUndeliveredMessagesPushRequestedUseCase(
	privateMessageRepo chatDomain.PrivateMessageRepository,
	privateMessageNotifier chatDomain.PrivateMessageNotifier,
	roomMessageRepo chatDomain.RoomMessageRepository,
	roomMessageNotifier chatDomain.RoomMessageNotifier,
) (chatApp.UndeliveredMessagesPushRequestedUseCase, error) {
	return chatApp.NewUndeliveredMessagesPushRequestedUseCase(
		privateMessageRepo,
		privateMessageNotifier,
		roomMessageRepo,
		roomMessageNotifier,
	)
}

func provideChatConfirmPrivateMessagesUseCase(
	messageRepo chatDomain.PrivateMessageRepository,
) (chatApp.ConfirmPrivateMessagesUseCase, error) {
	return chatApp.NewConfirmPrivateMessagesUseCase(
		messageRepo,
	)
}

func provideChatConfirmRoomMessagesUseCase(
	messageRepo chatDomain.RoomMessageRepository,
) (chatApp.ConfirmRoomMessagesUseCase, error) {
	return chatApp.NewConfirmRoomMessagesUseCase(
		messageRepo,
	)
}

func provideNotificationWelcomeEmailNotificationRequestedUseCase(
	emailNotifier notificationDomain.WelcomeEmailNotifier,
) (notificationApp.WelcomeEmailNotificationRequestedUseCase, error) {
	return notificationApp.NewWelcomeEmailNotificationRequestedUseCase(
		emailNotifier,
	)
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
func provideNotificationUndeliveredMessagesNotificationRequestedUseCase(
	systemMessageRepo notificationDomain.SystemMessageRepository,
	systemMessageNotifier notificationDomain.SystemMessageNotifier,
) (notificationApp.UndeliveredMessagesNotificationRequestedUseCase, error) {
	return notificationApp.NewUndeliveredMessagesNotificationRequestedUseCase(
		systemMessageRepo,
		systemMessageNotifier,
	)
}
func provideFriendshipUserCreatedUseCase(
	userRepo friendshipDomain.UserRepository,
) (friendshipApp.UserCreatedUseCase, error) {
	return friendshipApp.NewUserCreatedUseCase(
		userRepo,
	)
}
func provideFriendshipFriendRequestAgreedUseCase(
	requestRepo friendshipDomain.FriendRequestRepository,
	friendshipRepo friendshipDomain.FriendshipRepository,
	friendshipIDGen friendshipDomain.FriendshipIDGenerator,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendRequestAgreedUseCase, error) {
	return friendshipApp.NewFriendRequestAgreedUseCase(
		requestRepo,
		friendshipRepo,
		friendshipIDGen,
		eventIDGen,
	)
}
func provideFriendshipFriendRequestCreatedUseCase(
	friendRequestRepo friendshipDomain.FriendRequestRepository,
	eventRepo event.Repository,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendRequestCreatedUseCase, error) {
	return friendshipApp.NewFriendRequestCreatedUseCase(
		friendRequestRepo,
		eventRepo,
		eventIDGen,
	)
}
func provideFriendshipFriendshipCreatedUseCase(
	friendshipRepo friendshipDomain.FriendshipRepository,
	eventRepo event.Repository,
	eventIDGen event.IDGenerator,
) (friendshipApp.FriendshipCreatedUseCase, error) {
	return friendshipApp.NewFriendshipCreatedUseCase(
		friendshipRepo,
		eventRepo,
		eventIDGen,
	)
}
