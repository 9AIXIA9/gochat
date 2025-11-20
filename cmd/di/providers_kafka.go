package di

import (
	"context"
	"gochat/config"
	"gochat/internal/application/usecase"
	authorizationUsecase "gochat/internal/authorization/application"
	authorizationDomain "gochat/internal/authorization/domain"
	authorizationKafka "gochat/internal/authorization/port/kafka"
	chatUsecase "gochat/internal/chat/application/usecase"
	chatDomain "gochat/internal/chat/domain"
	chatKafka "gochat/internal/chat/port/kafka"
	"gochat/internal/delivery/kafka"
	kafkautil "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/infrastructure/websocket"
	notificationUsecase "gochat/internal/notification/application/usecase"
	notificationDomain "gochat/internal/notification/domain"
	notificationKafka "gochat/internal/notification/port/kafka"
	"gochat/internal/shared/event"
	"gochat/internal/social/application"
	socialDomain "gochat/internal/social/domain"
	socialKafka "gochat/internal/social/port/kafka"

	"github.com/google/wire"
	"go.uber.org/zap"
)

type (
	kafkaTopicEnsured    bool
	kafkaTopicSubscribed bool
)

var KafkaSet = wire.NewSet(
	provideTopicsEnsured,
	provideKafkaPublisher,
	provideKafkaSubscriber,
	provideKafkaTopicsSubscribed,
)

func provideKafkaPublisher(appConfig *config.App, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) (*kafkautil.EventPublisher, error) {
	return kafkautil.NewEventPublisher(appConfig.Kafka, eventRepo, metrics)
}

func provideKafkaSubscriber(appConfig *config.App, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) (*kafkautil.EventSubscriber, error) {
	return kafkautil.NewEventSubscriber(appConfig.Kafka, eventRepo, metrics)
}

func provideKafkaTopicsSubscribed(
	ensured kafkaTopicEnsured,
	subscriber *kafkautil.EventSubscriber,
	emailAvailable emailServiceAvailable,
	//websocket
	userSessionStartedUseCase usecase.UserSessionStartedUseCase,
	// auth
	authUserCreated authorizationUsecase.UserCreatedUseCase,
	// social
	socialUserCreated application.UserCreatedUseCase,
	socialRoomCreated application.RoomCreatedUseCase,
	socialRoomJoined application.RoomJoinedUseCase,
	socialRoomLeft application.RoomLeftUseCase,
	// chat
	chatUserCreated chatUsecase.UserCreatedUseCase,
	chatRoomCreated chatUsecase.RoomCreatedUseCase,
	chatRoomJoined chatUsecase.RoomJoinedUseCase,
	chatRoomLeft chatUsecase.RoomLeftUseCase,
	chatPrivateMessageCreated chatUsecase.PrivateMessageCreatedUseCase,
	chatRoomMessageCreated chatUsecase.RoomMessageCreatedUseCase,
	// notification
	notificationWelcomeEmailNotificationRequested notificationUsecase.WelcomeEmailNotificationRequestedUseCase,
	notificationMessageNotificationRequested notificationUsecase.MessageNotificationRequestedUseCase,
	notificationUndeliveredMessageNotificationRequested notificationUsecase.UndeliveredMessageNotificationRequestedUseCase,
) kafkaTopicSubscribed {
	if !ensured {
		zap.L().Warn("Kafka topics are not ensured, skipping subscription")
		return false
	}
	//websocket
	subscriber.Subscribe(websocket.TopicUserSessionStarted, kafka.NewUserSessionStartedEventHandler(userSessionStartedUseCase))
	// Authorization
	subscriber.Subscribe(authorizationDomain.TopicUserCreated, authorizationKafka.NewUserCreatedEventHandler(authUserCreated))
	// Social
	subscriber.Subscribe(socialDomain.TopicUserCreated, socialKafka.NewUserCreatedEventHandler(socialUserCreated))
	subscriber.Subscribe(socialDomain.TopicRoomCreated, socialKafka.NewRoomCreatedEventHandler(socialRoomCreated))
	subscriber.Subscribe(socialDomain.TopicRoomJoined, socialKafka.NewRoomJoinedEventHandler(socialRoomJoined))
	subscriber.Subscribe(socialDomain.TopicRoomLeft, socialKafka.NewRoomLeftEventHandler(socialRoomLeft))
	// Chat
	subscriber.Subscribe(chatDomain.TopicUserCreated, chatKafka.NewUserCreatedEventHandler(chatUserCreated))
	subscriber.Subscribe(chatDomain.TopicRoomCreated, chatKafka.NewRoomCreatedEventHandler(chatRoomCreated))
	subscriber.Subscribe(chatDomain.TopicRoomJoined, chatKafka.NewRoomJoinedEventHandler(chatRoomJoined))
	subscriber.Subscribe(chatDomain.TopicRoomLeft, chatKafka.NewRoomLeftEventHandler(chatRoomLeft))
	subscriber.Subscribe(chatDomain.TopicPrivateMessageCreated, chatKafka.NewPrivateMessageCreatedEventHandler(chatPrivateMessageCreated))
	subscriber.Subscribe(chatDomain.TopicRoomMessageCreated, chatKafka.NewRoomMessageCreatedEventHandler(chatRoomMessageCreated))
	// Notification
	if emailAvailable {
		zap.L().Info("Subscribing to WelcomeEmailNotificationRequested topic as email dialer is connected")
		subscriber.Subscribe(notificationDomain.TopicWelcomeEmailNotificationRequested, notificationKafka.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested))
	} else {
		zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
	}
	subscriber.Subscribe(notificationDomain.TopicMessageNotificationRequested, notificationKafka.NewMessageNotificationRequestedEventHandler(notificationMessageNotificationRequested))
	subscriber.Subscribe(notificationDomain.TopicUndeliveredMessageNotificationRequested, notificationKafka.NewUndeliveredMessageNotificationRequestedEventHandler(notificationUndeliveredMessageNotificationRequested))
	return true
}

func provideTopicsEnsured(appConfig *config.App) kafkaTopicEnsured {
	if err := kafkautil.EnsureTopics(
		context.Background(),
		appConfig.Kafka,
		[]event.Topic{
			// websocket
			websocket.TopicUserSessionStarted,
			// authorization
			authorizationDomain.TopicUserCreated,
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
			notificationDomain.TopicMessageNotificationRequested,
			notificationDomain.TopicUndeliveredMessageNotificationRequested,
		},
		1,
		1,
	); err != nil {
		zap.L().Warn("Failed to ensure Kafka topics", zap.Error(err))
		return false
	}
	return true
}
