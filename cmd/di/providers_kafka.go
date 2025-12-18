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
	"gochat/internal/delivery/kafka/handler"
	"gochat/internal/delivery/kafka/middleware"
	friendshipApp "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipEvent "gochat/internal/friendship/port/event"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	notificationEvent "gochat/internal/notification/port/event"
	profileApp "gochat/internal/profile/application"
	profileDomain "gochat/internal/profile/domain"
	profileEvent "gochat/internal/profile/port/event"
	roomshipApp "gochat/internal/roomship/application"
	roomshipDomain "gochat/internal/roomship/domain"
	roomshipEvent "gochat/internal/roomship/port/event"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
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
	provideKafkaProducer,
)

func provideKafkaConsumer(
	appConfig *config.App,
	router *kafkaInfra.Router,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
) (*kafkaInfra.Consumer, error) {
	consumer, err := kafkaInfra.NewConsumer(appConfig.Kafka, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer.SetErrorHandler(
		handler.NewLoggerErrorHandler(),
		middleware.NewRetryErrorMiddleware(reproducer),
		middleware.NewDeadLetterErrorMiddleware(
			eventRepo,
		),
	)

	return consumer, nil
}

func provideKafkaProducer(appConf *config.App) (*ckafka.Producer, error) {
	producer, err := kafkaInfra.NewProducer(appConf.Kafka)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	return producer, nil
}

// TODO 各上下文 分开订阅 添加中间件
func provideKafkaRouter(
	ensured kafkaTopicEnsured,
	emailAvailable emailServiceAvailable,
	appConfig *config.App,
	// auth
	authUserCreated authApp.UserCreatedUseCase,
	// profile
	profileUserCreated profileApp.UserCreatedUseCase,
	profileRoomCreated profileApp.RoomCreatedUseCase,
	profileRoomshipCreated profileApp.RoomshipCreatedUseCase,
	// chat
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomshipCreated chatApp.RoomshipCreatedUseCase,
	chatFriendshipCreated chatApp.FriendshipCreatedUseCase,
	chatUndeliveredMessagesPushRequested chatApp.UndeliveredMessagesPushRequestedUseCase,
	// notification
	notificationWelcomeEmailNotificationRequested notificationApp.WelcomeEmailNotificationRequestedUseCase,
	notificationSystemMessageNotificationRequested notificationApp.SystemMessageNotificationRequestedUseCase,
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
		router.EventHandle(authDomain.TopicUserCreated, authEvent.NewUserCreatedEventHandler(authUserCreated))
	}

	// Profile
	{
		router.EventHandle(profileDomain.TopicUserCreated, profileEvent.NewUserCreatedEventHandler(profileUserCreated))
		router.EventHandle(profileDomain.TopicRoomCreated, profileEvent.NewRoomCreatedEventHandler(profileRoomCreated))
		router.EventHandle(profileDomain.TopicRoomshipCreated, profileEvent.NewRoomshipCreatedEventHandler(profileRoomshipCreated))
	}

	// Chat
	{
		router.EventHandle(chatDomain.TopicUserCreated, chatEvent.NewUserCreatedEventHandler(chatUserCreated))
		router.EventHandle(chatDomain.TopicRoomCreated, chatEvent.NewRoomCreatedEventHandler(chatRoomCreated))
		router.EventHandle(chatDomain.TopicRoomshipCreated, chatEvent.NewRoomshipCreatedEventHandler(chatRoomshipCreated))
		router.EventHandle(chatDomain.TopicFriendshipCreated, chatEvent.NewFriendshipCreatedEventHandler(chatFriendshipCreated))
		router.EventHandle(chatDomain.TopicUndeliveredMessagesPushRequested, chatEvent.NewUndeliveredMessagesPushRequestedEventHandler(chatUndeliveredMessagesPushRequested))
	}
	// Notification
	{
		if emailAvailable {
			zap.L().Info("Subscribing to WelcomeEmailNotificationRequested topic as email dialer is connected")
			router.EventHandle(notificationDomain.TopicWelcomeEmailNotificationRequested, notificationEvent.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested))
		} else {
			zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
		}
		router.EventHandle(notificationDomain.TopicSystemMessageNotificationRequested, notificationEvent.NewSystemMessageNotificationRequestedEventHandler(notificationSystemMessageNotificationRequested))
		router.EventHandle(notificationDomain.TopicUndeliveredMessagesNotificationRequested, notificationEvent.NewUndeliveredMessagesNotificationRequestedEventHandler(notificationUndeliveredMessagesRequested))
	}

	// Roomship
	{
		router.EventHandle(roomshipDomain.TopicUserCreated, roomshipEvent.NewUserCreatedEventHandler(roomshipUserCreated))
		router.EventHandle(roomshipDomain.TopicRoomCreated, roomshipEvent.NewRoomCreatedEventHandler(roomshipRoomCreated))
		router.EventHandle(roomshipDomain.TopicMemberRequestCreated, roomshipEvent.NewMemberRequestCreatedEventHandler(roomshipMemberRequestCreated))
		router.EventHandle(roomshipDomain.TopicMemberRequestAgreed, roomshipEvent.NewMemberRequestAgreedEventHandler(roomshipMemberRequestAgreed))
		router.EventHandle(roomshipDomain.TopicRoomshipCreated, roomshipEvent.NewRoomshipCreatedEventHandler(roomshipRoomshipCreated))
	}

	// Friendship
	{
		router.EventHandle(friendshipDomain.TopicUserCreated, friendshipEvent.NewUserCreatedEventHandler(friendshipUserCreated))
		router.EventHandle(friendshipDomain.TopicFriendRequestAgreed, friendshipEvent.NewFriendRequestAgreedEventHandler(friendshipFriendRequestAgreed))
		router.EventHandle(friendshipDomain.TopicFriendRequestCreated, friendshipEvent.NewFriendRequestCreatedEventHandler(friendshipFriendRequestCreated))
		router.EventHandle(friendshipDomain.TopicFriendshipCreated, friendshipEvent.NewFriendshipCreatedEventHandler(friendshipFriendshipCreated))
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
			// profile
			profileDomain.TopicUserCreated,
			profileDomain.TopicRoomCreated,
			profileDomain.TopicRoomshipCreated,
			// roomship
			roomshipDomain.TopicUserCreated,
			roomshipDomain.TopicRoomshipCreated,
			roomshipDomain.TopicMemberRequestAgreed,
			roomshipDomain.TopicMemberRequestCreated,
			roomshipDomain.TopicRoomCreated,
			// chat
			chatDomain.TopicUserCreated,
			chatDomain.TopicRoomCreated,
			chatDomain.TopicFriendshipCreated,
			chatDomain.TopicRoomshipCreated,
			chatDomain.TopicUndeliveredMessagesPushRequested,
			// notification
			notificationDomain.TopicWelcomeEmailNotificationRequested,
			notificationDomain.TopicSystemMessageNotificationRequested,
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
