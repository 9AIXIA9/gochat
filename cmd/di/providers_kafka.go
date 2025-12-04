package di

import (
	"context"
	"fmt"
	"gochat/config"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authEvent "gochat/internal/authorization/port/event"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	chatEvent "gochat/internal/chat/port/event"
	"gochat/internal/delivery/kafka/middleware"
	friendshipApp "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipEvent "gochat/internal/friendship/port/event"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	notificationEvent "gochat/internal/notification/port/event"
	roomshipApp "gochat/internal/roomship/application"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipEvent "gochat/internal/roomship/port/event"
	"gochat/internal/shared/event"

	"github.com/google/wire"
	"go.uber.org/zap"
)

type (
	kafkaTopicEnsured bool
)

var KafkaSet = wire.NewSet(
	provideTopicsEnsured,
	provideKafkaRouter,
	provideKafkaConsumer,
)

func provideKafkaConsumer(
	appConfig *config.App,
	router *kafkaInfra.Router,
	creator event.DeadLetterCreator,
) (*kafkaInfra.Consumer, error) {
	consumer, err := kafkaInfra.NewConsumer(appConfig.Kafka, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	handler, err := kafkaInfra.NewErrorHandlerWithDeadLetterAndRetry(appConfig.Kafka, creator)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka error handler: %w", err)
	}

	consumer.SetErrorHandler(handler)
	return consumer, nil
}

// TODO 各上下文 分开订阅 添加中间件
func provideKafkaRouter(
	ensured kafkaTopicEnsured,
	emailAvailable emailServiceAvailable,
	appConfig *config.App,
	// auth
	authUserCreated authApp.UserCreatedUseCase,
	// chat
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomJoined chatApp.RoomJoinedUseCase,
	chatRoomLeft chatApp.RoomLeftUseCase,
	chatPrivateMessageCreated chatApp.PrivateMessageCreatedUseCase,
	chatRoomMessageCreated chatApp.RoomMessageCreatedUseCase,
	// notification
	notificationWelcomeEmailNotificationRequested notificationApp.WelcomeEmailNotificationRequestedUseCase,
	notificationFriendRequestCreatedNotificationRequested notificationApp.FriendRequestCreatedNotificationRequestedUseCase,
	notificationFriendshipCreatedNotificationRequested notificationApp.FriendshipCreatedNotificationRequestedUseCase,
	notificationPrivateMessageNotificationRequested notificationApp.PrivateMessageNotificationRequestedUseCase,
	notificationRoomMessageNotificationRequested notificationApp.RoomMessageNotificationRequestedUseCase,
	notificationUndeliveredMessagesRequested notificationApp.UndeliveredMessagesNotificationRequestedUseCase,
	// roomship
	roomshipUserCreated roomshipApp.UserCreatedUseCase,
	roomshipRoomCreated roomshipApp.RoomCreatedUseCase,
	roomshipMemberRequestAgreed roomshipApp.MemberRequestAgreedUseCase,
	roomshipMemberRequestCreated roomshipApp.MemberRequestCreatedUseCase,
	roomshipRoomshipCreated roomshipApp.RoomshipCreatedUseCase,
	// friendship
	friendshipUserCreated friendshipApp.UserCreatedUseCase,
	friendshipFriendRequestAgreed friendshipApp.FriendRequestAgreedUseCase,
	friendshipFriendRequestCreated friendshipApp.FriendRequestCreatedUseCase,
	friendshipFriendshipCreated friendshipApp.FriendshipCreatedUseCase,
) *kafkaInfra.Router {
	if !ensured {
		zap.L().Warn("Kafka topics are not ensured")
	}

	router := kafkaInfra.NewRouter()

	//TODO 可添加更多中间件 限流等
	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewTelemetryMiddleware(appConfig.Name),
	)

	// Authorization
	{
		router.Handle(authDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(authEvent.NewUserCreatedEventHandler(authUserCreated)))
	}

	// Chat
	{
		router.Handle(chatDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(chatEvent.NewUserCreatedEventHandler(chatUserCreated)))
		router.Handle(chatDomain.TopicRoomCreated, kafkaInfra.WrapEventHandler(chatEvent.NewRoomCreatedEventHandler(chatRoomCreated)))
		router.Handle(chatDomain.TopicRoomJoined, kafkaInfra.WrapEventHandler(chatEvent.NewRoomJoinedEventHandler(chatRoomJoined)))
		router.Handle(chatDomain.TopicRoomLeft, kafkaInfra.WrapEventHandler(chatEvent.NewRoomLeftEventHandler(chatRoomLeft)))
		router.Handle(chatDomain.TopicPrivateMessageCreated, kafkaInfra.WrapEventHandler(chatEvent.NewPrivateMessageCreatedEventHandler(chatPrivateMessageCreated)))
		router.Handle(chatDomain.TopicRoomMessageCreated, kafkaInfra.WrapEventHandler(chatEvent.NewRoomMessageCreatedEventHandler(chatRoomMessageCreated)))
	}
	// Notification
	{
		if emailAvailable {
			zap.L().Info("Subscribing to WelcomeEmailNotificationRequested topic as email dialer is connected")
			router.Handle(notificationDomain.TopicWelcomeEmailNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested)))
		} else {
			zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
		}
		router.Handle(notificationDomain.TopicFriendshipCreatedNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewFriendshipCreatedNotificationRequestedEventHandler(notificationFriendshipCreatedNotificationRequested)))
		router.Handle(notificationDomain.TopicFriendRequestCreatedNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewFriendRequestCreatedNotificationRequestedEventHandler(notificationFriendRequestCreatedNotificationRequested)))
		router.Handle(notificationDomain.TopicPrivateMessageNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewPrivateMessageNotificationRequestedEventHandler(notificationPrivateMessageNotificationRequested)))
		router.Handle(notificationDomain.TopicRoomMessageNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewRoomMessageNotificationRequestedEventHandler(notificationRoomMessageNotificationRequested)))
		router.Handle(notificationDomain.TopicUndeliveredMessagesNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewUndeliveredMessagesNotificationRequestedEventHandler(notificationUndeliveredMessagesRequested)))
	}

	// Roomship
	{
		router.Handle(roomshipDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(roomshipEvent.NewUserCreatedEventHandler(roomshipUserCreated)))
		router.Handle(roomshipDomain.TopicRoomCreated, kafkaInfra.WrapEventHandler(roomshipEvent.NewRoomCreatedEventHandler(roomshipRoomCreated)))
		router.Handle(roomshipDomain.TopicMemberRequestCreated, kafkaInfra.WrapEventHandler(roomshipEvent.NewMemberRequestCreatedEventHandler(roomshipMemberRequestCreated)))
		router.Handle(roomshipDomain.TopicMemberRequestAgreed, kafkaInfra.WrapEventHandler(roomshipEvent.NewMemberRequestAgreedEventHandler(roomshipMemberRequestAgreed)))
		router.Handle(roomshipDomain.TopicRoomshipCreated, kafkaInfra.WrapEventHandler(roomshipEvent.NewRoomshipCreatedEventHandler(roomshipRoomshipCreated)))
	}

	// Friendship
	{
		router.Handle(friendshipDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(friendshipEvent.NewUserCreatedEventHandler(friendshipUserCreated)))
		router.Handle(friendshipDomain.TopicFriendRequestAgreed, kafkaInfra.WrapEventHandler(friendshipEvent.NewFriendRequestAgreedEventHandler(friendshipFriendRequestAgreed)))
		router.Handle(friendshipDomain.TopicFriendRequestCreated, kafkaInfra.WrapEventHandler(friendshipEvent.NewFriendRequestCreatedEventHandler(friendshipFriendRequestCreated)))
		router.Handle(friendshipDomain.TopicFriendshipCreated, kafkaInfra.WrapEventHandler(friendshipEvent.NewFriendshipCreatedEventHandler(friendshipFriendshipCreated)))
	}

	return router
}

func provideTopicsEnsured(appConfig *config.App) kafkaTopicEnsured {
	if err := kafkaInfra.EnsureTopics(
		context.Background(),
		appConfig.Kafka,
		[]event.Topic{
			// authorization
			authDomain.TopicUserCreated,
			// roomship
			roomshipDomain.TopicUserCreated,
			roomshipDomain.TopicRoomshipCreated,
			roomshipDomain.TopicMemberRequestAgreed,
			roomshipDomain.TopicMemberRequestCreated,
			roomshipDomain.TopicRoomCreated,
			// chat
			chatDomain.TopicUserCreated,
			chatDomain.TopicRoomCreated,
			chatDomain.TopicRoomJoined,
			chatDomain.TopicRoomLeft,
			chatDomain.TopicPrivateMessageCreated,
			chatDomain.TopicRoomMessageCreated,
			// notification
			notificationDomain.TopicWelcomeEmailNotificationRequested,
			notificationDomain.TopicFriendshipCreatedNotificationRequested,
			notificationDomain.TopicFriendRequestCreatedNotificationRequested,
			notificationDomain.TopicPrivateMessageNotificationRequested,
			notificationDomain.TopicRoomMessageNotificationRequested,
			notificationDomain.TopicUndeliveredMessagesNotificationRequested,
			// friendship
			friendshipDomain.TopicUserCreated,
			friendshipDomain.TopicFriendRequestCreated,
			friendshipDomain.TopicFriendRequestAgreed,
			friendshipDomain.TopicFriendshipCreated,
		},
		1,
		1,
	); err != nil {
		zap.L().Warn("Failed to ensure Kafka topics", zap.Error(err))
		return false
	}
	return true
}
