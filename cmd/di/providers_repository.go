package di

import (
	authDomain "gochat/internal/authorization/domain"
	authConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authModel "gochat/internal/authorization/infrastructure/persistence/model"
	authRepo "gochat/internal/authorization/infrastructure/persistence/repository"
	chatDomain "gochat/internal/chat/domain"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepo "gochat/internal/chat/infrastructure/persistence/repository"
	gormInfra "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	notificationDomain "gochat/internal/notification/domain"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepo "gochat/internal/notification/infrastructure/persistence/repository"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialDomain "gochat/internal/social/domain"
	socialModel "gochat/internal/social/infrastructure/persistence/model"
	socialRepo "gochat/internal/social/infrastructure/persistence/repository"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type databaseMigrated bool

var RepoSet = wire.NewSet(
	wire.Bind(new(kernel.UnitOfWork), new(*gormInfra.UnitOfWork)),
	// Bind repositories to their interfaces
	wire.Bind(new(event.Repository), new(*repository.EventRepository)),

	wire.Bind(new(authDomain.UserRepository), new(*authRepo.UserRepository)),
	wire.Bind(new(authDomain.RefreshTokenRepository), new(*authRepo.RefreshTokenRepository)),

	wire.Bind(new(socialDomain.UserRepository), new(*socialRepo.UserRepository)),
	wire.Bind(new(socialDomain.RoomRepository), new(*socialRepo.RoomRepository)),

	wire.Bind(new(chatDomain.UserRepository), new(*chatRepo.UserRepository)),
	wire.Bind(new(chatDomain.RoomRepository), new(*chatRepo.RoomRepository)),
	wire.Bind(new(chatDomain.PrivateMessageRepository), new(*chatRepo.PrivateMessageRepository)),
	wire.Bind(new(chatDomain.RoomMessageRepository), new(*chatRepo.RoomMessageRepository)),

	wire.Bind(new(notificationDomain.PrivateMessageRepository), new(*notificationRepo.PrivateMessageRepository)),
	wire.Bind(new(notificationDomain.RoomMessageRepository), new(*notificationRepo.RoomMessageRepository)),

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
)

func provideDatabaseMigrated(mysql *gorm.DB) databaseMigrated {
	if err := gormInfra.AutoMigrate(
		mysql,
		&authModel.User{},
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

func provideUnitOfWork(mysql *gorm.DB, metrics *prometheus.Metrics) *gormInfra.UnitOfWork {
	u := gormInfra.NewUnitOfWork(mysql)
	u.SetMetrics(metrics)
	return u
}

func provideEventRepository(unitOfWork *gormInfra.UnitOfWork) *repository.EventRepository {
	return repository.NewEventRepository(unitOfWork)
}
func provideAuthorizationUserRepository(db *gorm.DB, eventRepo event.Repository) *authRepo.UserRepository {
	return authRepo.NewUserRepository(db, eventRepo)
}
func provideAuthorizationRefreshTokenRepository(redisClient *redis.Client) *authRepo.RefreshTokenRepository {
	return authRepo.NewRefreshTokenRepository(redisClient, &authConverter.RefreshTokenConverter{})
}
func provideChatUserRepository(db *gorm.DB) *chatRepo.UserRepository {
	return chatRepo.NewUserRepository(db)
}
func provideChatRoomRepository(db *gorm.DB) *chatRepo.RoomRepository {
	return chatRepo.NewRoomRepository(db)
}
func provideChatPrivateMessageRepository(db *gorm.DB, eventRepo event.Repository) *chatRepo.PrivateMessageRepository {
	return chatRepo.NewPrivateMessageRepository(db, eventRepo)
}
func provideChatRoomMessageRepository(db *gorm.DB, eventRepo event.Repository) *chatRepo.RoomMessageRepository {
	return chatRepo.NewRoomMessageRepository(db, eventRepo)
}
func provideSocialUserRepository(unitOfWork *gormInfra.UnitOfWork) *socialRepo.UserRepository {
	return socialRepo.NewUserRepository(unitOfWork)
}
func provideSocialRoomRepository(db *gorm.DB, eventRepo event.Repository) *socialRepo.RoomRepository {
	return socialRepo.NewRoomRepository(db, eventRepo)
}
func provideNotificationPrivateMessageRepository(unitOfWork *gormInfra.UnitOfWork) *notificationRepo.PrivateMessageRepository {
	return notificationRepo.NewPrivateMessageRepository(unitOfWork)
}
func provideNotificationRoomMessageRepository(unitOfWork *gormInfra.UnitOfWork) *notificationRepo.RoomMessageRepository {
	return notificationRepo.NewRoomMessageRepository(unitOfWork)
}
