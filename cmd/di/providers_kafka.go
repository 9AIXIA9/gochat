package di

import (
	"context"
	"gochat/config"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authEvent "gochat/internal/authorization/port/event"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	chatEvent "gochat/internal/chat/port/event"
	"gochat/internal/delivery/kafka/middleware"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	notificationEvent "gochat/internal/notification/port/event"
	"gochat/internal/shared/event"
	socialApp "gochat/internal/social/application"
	socialDomain "gochat/internal/social/domain"
	socialEvent "gochat/internal/social/port/event"

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

func provideKafkaConsumer(appConfig *config.App, router *kafkaInfra.Router) (*kafkaInfra.Consumer, error) {
	return kafkaInfra.NewConsumer(appConfig.Kafka, router)
}

// TODO 各上下文 分开订阅 添加中间件
func provideKafkaRouter(
	ensured kafkaTopicEnsured,
	emailAvailable emailServiceAvailable,
	appConfig *config.App,
	// auth
	authUserCreated authApp.UserCreatedUseCase,
	// social
	socialUserCreated socialApp.UserCreatedUseCase,
	socialRoomCreated socialApp.RoomCreatedUseCase,
	socialRoomJoined socialApp.RoomJoinedUseCase,
	socialRoomLeft socialApp.RoomLeftUseCase,
	// chat
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomJoined chatApp.RoomJoinedUseCase,
	chatRoomLeft chatApp.RoomLeftUseCase,
	chatPrivateMessageCreated chatApp.PrivateMessageCreatedUseCase,
	chatRoomMessageCreated chatApp.RoomMessageCreatedUseCase,
	// notification
	notificationWelcomeEmailNotificationRequested notificationApp.WelcomeEmailNotificationRequestedUseCase,
	notificationPrivateMessageNotificationRequested notificationApp.PrivateMessageNotificationRequestedUseCase,
	notificationRoomMessageNotificationRequested notificationApp.RoomMessageNotificationRequestedUseCase,
	notificationUndeliveredMessagesRequested notificationApp.UndeliveredMessagesNotificationRequestedUseCase,
) *kafkaInfra.Router {
	if !ensured {
		zap.L().Warn("Kafka topics are not ensured")
	}

	router := kafkaInfra.NewRouter()

	//TODO 可添加更多中间件 重试，死信，限流等
	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewTelemetryMiddleware(appConfig.Name),
	)

	// Authorization
	router.Handle(authDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(authEvent.NewUserCreatedEventHandler(authUserCreated)))
	// Social
	router.Handle(socialDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(socialEvent.NewUserCreatedEventHandler(socialUserCreated)))
	router.Handle(socialDomain.TopicRoomCreated, kafkaInfra.WrapEventHandler(socialEvent.NewRoomCreatedEventHandler(socialRoomCreated)))
	router.Handle(socialDomain.TopicRoomJoined, kafkaInfra.WrapEventHandler(socialEvent.NewRoomJoinedEventHandler(socialRoomJoined)))
	router.Handle(socialDomain.TopicRoomLeft, kafkaInfra.WrapEventHandler(socialEvent.NewRoomLeftEventHandler(socialRoomLeft)))
	// Chat
	router.Handle(chatDomain.TopicUserCreated, kafkaInfra.WrapEventHandler(chatEvent.NewUserCreatedEventHandler(chatUserCreated)))
	router.Handle(chatDomain.TopicRoomCreated, kafkaInfra.WrapEventHandler(chatEvent.NewRoomCreatedEventHandler(chatRoomCreated)))
	router.Handle(chatDomain.TopicRoomJoined, kafkaInfra.WrapEventHandler(chatEvent.NewRoomJoinedEventHandler(chatRoomJoined)))
	router.Handle(chatDomain.TopicRoomLeft, kafkaInfra.WrapEventHandler(chatEvent.NewRoomLeftEventHandler(chatRoomLeft)))
	router.Handle(chatDomain.TopicPrivateMessageCreated, kafkaInfra.WrapEventHandler(chatEvent.NewPrivateMessageCreatedEventHandler(chatPrivateMessageCreated)))
	router.Handle(chatDomain.TopicRoomMessageCreated, kafkaInfra.WrapEventHandler(chatEvent.NewRoomMessageCreatedEventHandler(chatRoomMessageCreated)))
	// Notification
	if emailAvailable {
		zap.L().Info("Subscribing to WelcomeEmailNotificationRequested topic as email dialer is connected")
		router.Handle(notificationDomain.TopicWelcomeEmailNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested)))
	} else {
		zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
	}
	router.Handle(notificationDomain.TopicPrivateMessageNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewPrivateMessageNotificationRequestedEventHandler(notificationPrivateMessageNotificationRequested)))
	router.Handle(notificationDomain.TopicRoomMessageNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewRoomMessageNotificationRequestedEventHandler(notificationRoomMessageNotificationRequested)))
	router.Handle(notificationDomain.TopicUndeliveredMessagesNotificationRequested, kafkaInfra.WrapEventHandler(notificationEvent.NewUndeliveredMessagesNotificationRequestedEventHandler(notificationUndeliveredMessagesRequested)))
	return router
}

func provideTopicsEnsured(appConfig *config.App) kafkaTopicEnsured {
	if err := kafkaInfra.EnsureTopics(
		context.Background(),
		appConfig.Kafka,
		[]event.Topic{
			// authorization
			authDomain.TopicUserCreated,
			// social
			socialDomain.TopicUserCreated,
			socialDomain.TopicRoomCreated,
			socialDomain.TopicRoomJoined,
			socialDomain.TopicRoomLeft,
			// chat
			chatDomain.TopicUserCreated,
			chatDomain.TopicRoomCreated,
			chatDomain.TopicRoomJoined,
			chatDomain.TopicRoomLeft,
			chatDomain.TopicPrivateMessageCreated,
			chatDomain.TopicRoomMessageCreated,
			// notification
			notificationDomain.TopicWelcomeEmailNotificationRequested,
			notificationDomain.TopicPrivateMessageNotificationRequested,
			notificationDomain.TopicRoomMessageNotificationRequested,
			notificationDomain.TopicUndeliveredMessagesNotificationRequested,
		},
		1,
		1,
	); err != nil {
		zap.L().Warn("Failed to ensure Kafka topics", zap.Error(err))
		return false
	}
	return true
}
