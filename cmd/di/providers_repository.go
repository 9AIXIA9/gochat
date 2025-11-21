package di

import (
	"gochat/internal/authorization/domain"
	authorizationConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authorizationModel "gochat/internal/authorization/infrastructure/persistence/model"
	authorizationRepository "gochat/internal/authorization/infrastructure/persistence/repository"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepository "gochat/internal/chat/infrastructure/persistence/repository"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepository "gochat/internal/notification/infrastructure/persistence/repository"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialApp "gochat/internal/social/domain"
	socialModel "gochat/internal/social/infrastructure/persistence/model"
	socialRepository "gochat/internal/social/infrastructure/persistence/repository"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type databaseMigrated bool

var RepoSet = wire.NewSet(
	provideDatabaseMigrated,
	provideUnitOfWork,
	provideEventRepository,
	provideAuthorizationUserRepository,
	provideAuthorizationRefreshTokenRepository,
	provideChatUserRepository,
	provideChatRoomRepository,
	provideChatPrivateMessageRepository,
	provideChatRoomMessageRepository,
	provideSocialUserRepository,
	provideSocialRoomRepository,
	provideNotificationPrivateMessageRepository,
	provideNotificationRoomMessageRepository,
	// Binds whose concrete types are provided in this set
	wire.Bind(new(kernel.UnitOfWork), new(*gormutils.UnitOfWork)),
	// Event store binds
	wire.Bind(new(event.UnpublishedEventSaver), new(*repository.EventRepository)),
	wire.Bind(new(event.UnpublishedEventsSaver), new(*repository.EventRepository)),
	// Authorization repository binds
	wire.Bind(new(domain.UserSaver), new(*authorizationRepository.UserRepository)),
	wire.Bind(new(domain.UserFinderByNumber), new(*authorizationRepository.UserRepository)),
	wire.Bind(new(domain.RefreshTokenSaver), new(*authorizationRepository.RefreshTokenRepository)),
	wire.Bind(new(domain.RefreshTokenFinder), new(*authorizationRepository.RefreshTokenRepository)),
	// Social repository binds
	wire.Bind(new(socialApp.RoomSaver), new(*socialRepository.RoomRepository)),
	wire.Bind(new(socialApp.RoomFinderByNumber), new(*socialRepository.RoomRepository)),
	wire.Bind(new(socialApp.RoomMemberSaver), new(*socialRepository.RoomRepository)),
	wire.Bind(new(socialApp.RoomMemberDeleter), new(*socialRepository.RoomRepository)),
)

func provideDatabaseMigrated(mysql *gorm.DB) databaseMigrated {
	if err := gormutils.AutoMigrate(
		mysql,
		&authorizationModel.User{},
		&notificationModel.PrivateMessage{},
		&notificationModel.RoomMessage{},
		&notificationModel.RoomMessageRecipient{},
		&notificationModel.RoomMessageState{},
		&chatModel.User{},
		&chatModel.Room{},
		&chatModel.PrivateMessage{},
		&chatModel.RoomMessage{},
		&socialModel.User{},
		&socialModel.Room{},
		&model.Event{},
		&model.DeadLetter{},
	); err != nil {
		zap.L().Warn("failed to migrate database", zap.Error(err))
		return false
	}
	return true
}

func provideUnitOfWork(mysql *gorm.DB, metrics *prometheus.Metrics) *gormutils.UnitOfWork {
	u := gormutils.NewUnitOfWork(mysql)
	u.SetMetrics(metrics)
	return u
}

func provideEventRepository(unitOfWork *gormutils.UnitOfWork) *repository.EventRepository {
	return repository.NewEventRepository(unitOfWork)
}
func provideAuthorizationUserRepository(unitOfWork *gormutils.UnitOfWork) *authorizationRepository.UserRepository {
	return authorizationRepository.NewUserRepository(unitOfWork)
}
func provideAuthorizationRefreshTokenRepository(redisClient *redis.Client) *authorizationRepository.RefreshTokenRepository {
	return authorizationRepository.NewRefreshTokenRepository(redisClient, &authorizationConverter.RefreshTokenConverter{})
}
func provideChatUserRepository(unitOfWork *gormutils.UnitOfWork) *chatRepository.UserRepository {
	return chatRepository.NewUserRepository(unitOfWork)
}
func provideChatRoomRepository(unitOfWork *gormutils.UnitOfWork) *chatRepository.RoomRepository {
	return chatRepository.NewRoomRepository(unitOfWork)
}
func provideChatPrivateMessageRepository(unitOfWork *gormutils.UnitOfWork) *chatRepository.PrivateMessageRepository {
	return chatRepository.NewPrivateMessageRepository(unitOfWork)
}
func provideChatRoomMessageRepository(unitOfWork *gormutils.UnitOfWork) *chatRepository.RoomMessageRepository {
	return chatRepository.NewRoomMessageRepository(unitOfWork)
}
func provideSocialUserRepository(unitOfWork *gormutils.UnitOfWork) *socialRepository.UserRepository {
	return socialRepository.NewUserRepository(unitOfWork)
}
func provideSocialRoomRepository(unitOfWork *gormutils.UnitOfWork) *socialRepository.RoomRepository {
	return socialRepository.NewRoomRepository(unitOfWork)
}
func provideNotificationPrivateMessageRepository(unitOfWork *gormutils.UnitOfWork) *notificationRepository.PrivateMessageRepository {
	return notificationRepository.NewPrivateMessageRepository(unitOfWork)
}
func provideNotificationRoomMessageRepository(unitOfWork *gormutils.UnitOfWork) *notificationRepository.RoomMessageRepository {
	return notificationRepository.NewRoomMessageRepository(unitOfWork)
}
