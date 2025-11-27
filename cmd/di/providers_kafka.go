package di

import (
	"context"
	"gochat/config"
	authApp "gochat/internal/authorization/application"
	authDomain "gochat/internal/authorization/domain"
	authKafka "gochat/internal/authorization/port/kafka"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	chatKafka "gochat/internal/chat/port/kafka"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/persistence/repository"
	"gochat/internal/infrastructure/prometheus"
	notificationApp "gochat/internal/notification/application"
	notificationDomain "gochat/internal/notification/domain"
	notificationKafka "gochat/internal/notification/port/kafka"
	"gochat/internal/shared/event"
	socialApp "gochat/internal/social/application"
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
	provideKafkaSubscriber,
	provideKafkaTopicsSubscribed,
)

func provideKafkaSubscriber(appConfig *config.App, eventRepo *repository.EventRepository, metrics *prometheus.Metrics) (*kafkaInfra.EventSubscriber, error) {
	return kafkaInfra.NewEventSubscriber(appConfig.Kafka, eventRepo, metrics)
}

func provideKafkaTopicsSubscribed(
	ensured kafkaTopicEnsured,
	subscriber *kafkaInfra.EventSubscriber,
	emailAvailable emailServiceAvailable,
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
) kafkaTopicSubscribed {
	if !ensured {
		zap.L().Warn("Kafka topics are not ensured, skipping subscription")
		return false
	}
	// Authorization
	subscriber.Subscribe(authDomain.TopicUserCreated, authKafka.NewUserCreatedEventHandler(authUserCreated))
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
	subscriber.Subscribe(notificationDomain.TopicPrivateMessageNotificationRequested, notificationKafka.NewPrivateMessageNotificationRequestedEventHandler(notificationPrivateMessageNotificationRequested))
	subscriber.Subscribe(notificationDomain.TopicRoomMessageNotificationRequested, notificationKafka.NewRoomMessageNotificationRequestedEventHandler(notificationRoomMessageNotificationRequested))
	subscriber.Subscribe(notificationDomain.TopicUndeliveredMessagesNotificationRequested, notificationKafka.NewUndeliveredMessagesNotificationRequestedEventHandler(notificationUndeliveredMessagesRequested))
	return true
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
