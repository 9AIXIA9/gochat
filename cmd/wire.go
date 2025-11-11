//go:build wireinject
// +build wireinject

package main

import (
	authorizationApp "gochat/internal/authorization/application"
	authorizationCrypto "gochat/internal/authorization/infrastructure/crypto"
	authorizationJWT "gochat/internal/authorization/infrastructure/jwt"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	authorizationSnowflake "gochat/internal/authorization/infrastructure/snowflake"
	authorizationUuid "gochat/internal/authorization/infrastructure/uuid"
	chatApp "gochat/internal/chat/application"
	chatUUID "gochat/internal/chat/infrastructure/uuid"
	"gochat/internal/infrastructure/bcrypt"
	infraRepository "gochat/internal/infrastructure/persistence/repository"
	infraUUID "gochat/internal/infrastructure/uuid"
	"gochat/internal/shared/event"
	socialApp "gochat/internal/social/application"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUUID "gochat/internal/social/infrastructure/uuid"

	"github.com/google/wire"
)

//go:generate go run github.com/google/wire/cmd/wire@latest ./...

// initializeDependencies builds the application Dependencies using Google Wire.
func initializeDependencies(configPath ConfigPath, envPath EnvPath) (*Dependencies, error) {
	wire.Build(
		// Base config & infra
		provideAppConfig,
		provideMysql,
		provideRedis,
		provideValidator,

		// Interface bindings to concrete providers returned by our providers
		wire.Bind(new(event.IDGenerator), new(*infraUUID.EventIDGenerator)),
		wire.Bind(new(authorizationApp.UserIDGenerator), new(*authorizationUuid.UserIDGenerator)),
		wire.Bind(new(authorizationApp.UserNumberGenerator), new(*authorizationSnowflake.UserNumberGenerator)),
		wire.Bind(new(authorizationApp.Encryptor), new(*bcrypt.Hasher)),
		wire.Bind(new(authorizationApp.Comparator), new(*bcrypt.Hasher)),
		wire.Bind(new(authorizationApp.UserSaver), new(*authorizationRepository.UserRepository)),
		wire.Bind(new(authorizationApp.UserFinderByNumber), new(*authorizationRepository.UserRepository)),
		wire.Bind(new(authorizationApp.UserLoggedInAtUpdater), new(*authorizationRepository.UserRepository)),
		wire.Bind(new(authorizationApp.RefreshTokenSaver), new(*authorizationRepository.RefreshTokenRepository)),
		wire.Bind(new(authorizationApp.RefreshTokenFinder), new(*authorizationRepository.RefreshTokenRepository)),
		wire.Bind(new(authorizationApp.AccessTokenGenerator), new(*authorizationJWT.AccessTokenManager)),
		wire.Bind(new(authorizationApp.AccessTokenParser), new(*authorizationJWT.AccessTokenManager)),
		wire.Bind(new(authorizationApp.RefreshTokenGenerator), new(*authorizationCrypto.RefreshTokenGenerator)),
		wire.Bind(new(chatApp.MessageIDGenerator), new(*chatUUID.MessageIDGenerator)),
		wire.Bind(new(event.UnpublishedSaver), new(*infraRepository.EventRepository)),
		wire.Bind(new(socialApp.RoomIDGenerator), new(*socialUUID.RoomIDGenerator)),
		wire.Bind(new(socialApp.RoomNumberGenerator), new(*socialSnowflake.RoomNumberGenerator)),
		wire.Bind(new(socialApp.Encryptor), new(*bcrypt.Hasher)),
		wire.Bind(new(socialApp.Comparator), new(*bcrypt.Hasher)),
		wire.Bind(new(socialApp.RoomSaver), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomFinder), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomJoiner), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomLeaver), new(*socialRepository.RoomRepository)),

		// Generators & managers (concrete providers)
		provideEventIDGenerator,
		provideAuthorizationUserIDGenerator,
		provideAuthorizationUserNumberGenerator,
		provideHasher,
		provideMessageIDGenerator,
		provideSocialRoomIDGenerator,
		provideSocialRoomNumberGenerator,
		provideAccessTokenManager,
		provideRefreshTokenGenerator,
		provideEmailNotifier,
		provideWebsocketManager,
		// Repositories
		provideEventRepository,
		provideAuthorizationUserRepository,
		provideAuthorizationRefreshTokenRepository,
		provideChatUserRepository,
		provideChatRoomRepository,
		provideChatMessageRepository,
		provideSocialUserRepository,
		provideSocialRoomRepository,
		provideNotificationUserRepository,
		provideNotificationRoomRepository,
		provideNotificationMessageRepository,
		// Kafka & Canal
		provideKafkaPublisher,
		provideKafkaSubscriber,
		provideOutboxConsumer,
		// HTTP UseCases
		provideSignUpUseCase,
		provideLoginUseCase,
		provideRefreshAccessTokenUseCase,
		provideParseAccessTokenUseCase,
		provideSendPrivateMessageUseCase,
		provideSendRoomMessageUseCase,
		provideCreateRoomUseCase,
		provideJoinRoomUseCase,
		provideLeaveRoomUseCase,
		provideNotificationUserConnectedUseCase,
		// Event UseCases (consumer side)
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
		provideNotificationUserCreatedUseCase,
		provideNotificationRoomCreatedUseCase,
		provideNotificationRoomJoinedUseCase,
		provideNotificationRoomLeftUseCase,
		provideNotificationPrivateMessageCreatedUseCase,
		provideNotificationRoomMessageCreatedUseCase,
		// Subscriptions (side-effect) must run before assembling final deps
		provideKafkaSubscriptions,
		// Final assembler
		BuildDependencies,
	)
	return &Dependencies{}, nil
}
