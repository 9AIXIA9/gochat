package di

import (
	authDomain "gochat/internal/authorization/domain"
	authConverter "gochat/internal/authorization/infrastructure/persistence/converter"
	authModel "gochat/internal/authorization/infrastructure/persistence/model"
	authRepo "gochat/internal/authorization/infrastructure/persistence/repository"
	chatDomain "gochat/internal/chat/domain"
	chatModel "gochat/internal/chat/infrastructure/persistence/model"
	chatRepo "gochat/internal/chat/infrastructure/persistence/repository"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipModel "gochat/internal/friendship/infrastructure/persistence/model"
	friendshipRepo "gochat/internal/friendship/infrastructure/persistence/repository"
	gormInfra "gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/persistence/model"
	"gochat/internal/infrastructure/persistence/repository"
	notificationDomain "gochat/internal/notification/domain"
	notificationModel "gochat/internal/notification/infrastructure/persistence/model"
	notificationRepo "gochat/internal/notification/infrastructure/persistence/repository"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipModel "gochat/internal/roomship/infrastructure/persistence/model"
	roomshipRepo "gochat/internal/roomship/infrastructure/persistence/repository"
	"gochat/internal/shared/event"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type databaseMigrated bool

var RepoSet = wire.NewSet(
	// Bind repositories to their interfaces
	wire.Bind(new(event.Repository), new(*repository.EventRepository)),

	wire.Bind(new(authDomain.UserRepository), new(*authRepo.UserRepository)),
	wire.Bind(new(authDomain.RefreshTokenRepository), new(*authRepo.RefreshTokenRepository)),

	wire.Bind(new(roomshipDomain.UserRepository), new(*roomshipRepo.UserRepository)),
	wire.Bind(new(roomshipDomain.RoomRepository), new(*roomshipRepo.RoomRepository)),

	wire.Bind(new(chatDomain.UserRepository), new(*chatRepo.UserRepository)),
	wire.Bind(new(chatDomain.RoomRepository), new(*chatRepo.RoomRepository)),
	wire.Bind(new(chatDomain.PrivateMessageRepository), new(*chatRepo.PrivateMessageRepository)),
	wire.Bind(new(chatDomain.RoomMessageRepository), new(*chatRepo.RoomMessageRepository)),

	wire.Bind(new(notificationDomain.PrivateMessageRepository), new(*notificationRepo.PrivateMessageRepository)),
	wire.Bind(new(notificationDomain.RoomMessageRepository), new(*notificationRepo.RoomMessageRepository)),

	wire.Bind(new(friendshipDomain.UserRepository), new(*friendshipRepo.UserRepository)),
	wire.Bind(new(friendshipDomain.FriendRequestRepository), new(*friendshipRepo.FriendRequestRepository)),

	wire.Bind(new(event.DeadLetterCreator), new(*repository.EventRepository)),

	provideDatabaseMigrated,
	provideEventRepository,
	provideAuthorizationUserRepository,
	provideAuthorizationRefreshTokenRepository,
	provideChatUserRepository,
	provideChatRoomRepository,
	provideChatPrivateMessageRepository,
	provideChatRoomMessageRepository,
	provideRoomshipUserRepository,
	provideRoomshipRoomRepository,
	provideNotificationPrivateMessageRepository,
	provideNotificationRoomMessageRepository,
	provideFriendshipUserRepository,
	provideFriendshipFriendRequestRepository,
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
		&roomshipModel.User{},
		&roomshipModel.Room{},
		&friendshipModel.User{},
		&friendshipModel.Friendship{},
		&friendshipModel.FriendRequest{},
		&model.Event{},
		&model.DeadLetter{},
	); err != nil {
		zap.L().Warn("failed to migrate database", zap.Error(err))
		return false
	}
	return true
}

func provideEventRepository(db *gorm.DB) *repository.EventRepository {
	return repository.NewEventRepository(db)
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
func provideRoomshipUserRepository(db *gorm.DB) *roomshipRepo.UserRepository {
	return roomshipRepo.NewUserRepository(db)
}
func provideRoomshipRoomRepository(db *gorm.DB, eventRepo event.Repository) *roomshipRepo.RoomRepository {
	return roomshipRepo.NewRoomRepository(db, eventRepo)
}
func provideNotificationPrivateMessageRepository(db *gorm.DB) *notificationRepo.PrivateMessageRepository {
	return notificationRepo.NewPrivateMessageRepository(db)
}
func provideNotificationRoomMessageRepository(db *gorm.DB) *notificationRepo.RoomMessageRepository {
	return notificationRepo.NewRoomMessageRepository(db)
}
func provideFriendshipUserRepository(db *gorm.DB, eventRepo event.Repository) *friendshipRepo.UserRepository {
	return friendshipRepo.NewUserRepository(db, eventRepo)
}
func provideFriendshipFriendRequestRepository(db *gorm.DB) *friendshipRepo.FriendRequestRepository {
	return friendshipRepo.NewFriendRequestRepository(db)
}
