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
	gormutils "gochat/internal/infrastructure/gorm"
	infraRepository "gochat/internal/infrastructure/persistence/repository"
	infraUUID "gochat/internal/infrastructure/uuid"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialApp "gochat/internal/social/application"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"
	socialSnowflake "gochat/internal/social/infrastructure/snowflake"
	socialUUID "gochat/internal/social/infrastructure/uuid"

	"github.com/google/wire"
)

//go:generate go run wire ./...

// initializeDependencies builds the application Dependencies using Google Wire.
func initializeDependencies(configPath ConfigPath, envPath EnvPath) (*Dependencies, error) {
	wire.Build(
		// Base config & infra
		provideAppConfig,
		provideMysql,
		provideRedis,
		provideValidator,
		provideMetrics,

		// Interface bindings to concrete providers returned by our providers
		wire.Bind(new(kernel.UnitOfWork), new(*gormutils.UnitOfWork)),
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
		wire.Bind(new(event.UnpublishedEventSaver), new(*infraRepository.EventRepository)),
		wire.Bind(new(event.UnpublishedEventsSaver), new(*infraRepository.EventRepository)),
		wire.Bind(new(socialApp.RoomIDGenerator), new(*socialUUID.RoomIDGenerator)),
		wire.Bind(new(socialApp.RoomNumberGenerator), new(*socialSnowflake.RoomNumberGenerator)),
		wire.Bind(new(socialApp.Encryptor), new(*bcrypt.Hasher)),
		wire.Bind(new(socialApp.Comparator), new(*bcrypt.Hasher)),
		wire.Bind(new(socialApp.RoomSaver), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomFinderByNumber), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomMemberSaver), new(*socialRepository.RoomRepository)),
		wire.Bind(new(socialApp.RoomMemberDeleter), new(*socialRepository.RoomRepository)),

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
		provideGomailDialer,
		provideEmailAvailable,
		provideEmailNotifier,
		provideMessageNotifier,
		// Repositories
		provideUnitOfWork,
		provideEventRepository,
		provideAuthorizationUserRepository,
		provideAuthorizationRefreshTokenRepository,
		provideChatUserRepository,
		provideChatRoomRepository,
		provideChatMessageRepository,
		provideSocialUserRepository,
		provideSocialRoomRepository,
		provideNotificationMessageRepository,
		// Kafka & Canal & Websocket
		provideKafkaTopics,
		provideKafkaPublisher,
		provideKafkaSubscriber,
		provideCanal,
		provideCanalOutboxConsumer,
		provideWebsocketUpgrader,
		provideWebsocketManager,
		provideWebsocketServer,
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
		provideNotificationUndeliveredMessageNotificationRequestedUseCase,
		provideNotificationMessageReadUseCase,
		// Event UseCases (consumer side)
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
		// HTTP
		provideHttpRouter,
		// Websocket
		provideWebsocketRouter,
		// Subscriptions (side-effect) must run before assembling final deps
		provideKafkaSubscriptions,
		// Final assembler
		BuildDependencies,
	)
	return &Dependencies{}, nil
}
