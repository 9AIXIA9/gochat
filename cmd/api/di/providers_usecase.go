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
)

var UseCaseKafkaSet = wire.NewSet(
	provideUnpublishedEventsCreatedCase,
	provideProfileRoomshipCreatedUseCase,
	provideProfileRoomCreatedUseCase,
	provideProfileUserCreatedUseCase,
	provideRoomshipUserCreatedUseCase,
	provideRoomshipMemberRequestAgreedUseCase,
	provideRoomshipRoomCreatedUseCase,
	provideChatUserCreatedUseCase,
	provideChatRoomCreatedUseCase,
	provideChatRoomshipCreatedUseCase,
	provideChatFriendshipCreatedUseCase,
	provideFriendshipUserCreatedUseCase,
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
	generator event.IDGenerator,
	messageIDGen kernel.MessageIDGenerator,
	messageRepo chatDomain.PrivateMessageRepository,
	friendshipRepo chatDomain.FriendshipRepository,
) (chatApp.SendPrivateMessageUseCase, error) {
	return chatApp.NewSendPrivateMessageUseCase(
		friendshipRepo,
		messageIDGen,
		messageRepo,
		generator,
	)
}
func provideSendRoomMessageUseCase(
	messageIDGen kernel.MessageIDGenerator,
	generator event.IDGenerator,
	roomshipRepo chatDomain.RoomshipRepository,
	messageRepo chatDomain.RoomMessageRepository,
) (chatApp.SendRoomMessageUseCase, error) {
	return chatApp.NewSendRoomMessageUseCase(
		messageIDGen,
		generator,
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
	operationIDGenerator kernel.OperationIDGenerator,
	comparator roomshipDomain.Comparator,
) (roomshipApp.SendMemberRequestUseCase, error) {
	return roomshipApp.NewSendMemberRequestUseCase(
		requestRepo,
		roomshipRepo,
		roomRepo,
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
	operationIDGenerator kernel.OperationIDGenerator,
) (friendshipApp.SendFriendRequestUseCase, error) {
	return friendshipApp.NewSendFriendRequestUseCase(
		friendshipRepo,
		requestRepo,
		requestRepo,
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

// -------------------- Event UseCases (websocket side) --------------------
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
) (roomshipApp.RoomCreatedUseCase, error) {
	return roomshipApp.NewRoomCreatedUseCase(
		roomshipIDGenerator,
		idGenerator,
		roomRepo,
		roomshipRepo,
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
