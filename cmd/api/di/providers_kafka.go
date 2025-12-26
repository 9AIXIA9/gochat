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

var KafkaSet = wire.NewSet(
	provideKafkaProducer,
	provideAuthEventConsumer,
	provideProfileEventConsumer,
	provideChatEventConsumer,
	provideNotificationEventConsumer,
	provideRoomshipEventConsumer,
	provideFriendshipEventConsumer,
)

// Distinct wrapper types for per-context consumers to avoid Wire ambiguity
// Underlying type is *kafkaInfra.Consumer, but each is a separate named type
// so Wire can differentiate providers and parameters.
type (
	AuthKafkaConsumer         *kafkaInfra.Consumer
	ProfileKafkaConsumer      *kafkaInfra.Consumer
	ChatKafkaConsumer         *kafkaInfra.Consumer
	NotificationKafkaConsumer *kafkaInfra.Consumer
	RoomshipKafkaConsumer     *kafkaInfra.Consumer
	FriendshipKafkaConsumer   *kafkaInfra.Consumer
)

func provideKafkaProducer(
	appConf *config.App,
) (*ckafka.Producer, error) {
	producer, err := kafkaInfra.NewProducer(appConf.Kafka)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	return producer, nil
}

// common builder to reduce duplication across contexts
func buildKafkaConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	contextName string,
	topics []event.Topic,
	register func(r *kafkaInfra.Router),
) (*kafkaInfra.Consumer, error) {
	if appConfig.Env == "dev" || appConfig.Env == "development" {
		if err := kafkaInfra.EnsureTopics(
			context.Background(),
			appConfig.Kafka,
			topics,
			1,
			1,
		); err != nil {
			zap.L().Warn(
				"Failed to ensure Kafka topics",
				zap.Strings("topics", func() []string {
					var ts []string
					for _, t := range topics {
						ts = append(ts, t.String())
					}
					return ts
				}()),
				zap.Error(err),
			)
		}
	}

	router := kafkaInfra.NewRouter()
	router.Use(
		middleware.NewLoggerMiddleware(),
		middleware.NewRecoverMiddleware(),
		middleware.NewTraceMiddleware(appConfig.Name+"."+contextName+".kafka_consumer"),
	)

	if appConfig.Breaker != nil {
		conf := *appConfig.Breaker
		conf.Name = appConfig.Name + "_" + contextName + "_kafka_consumer_circuit_breaker"
		router.Use(middleware.NewCircuitBreakMiddleware(&conf))
	}
	register(router)

	consumer, err := kafkaInfra.NewConsumer(appConfig.Kafka, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer.SetErrorHandler(
		handler.NewLoggerErrorHandler(),
		middleware.NewRetryErrorMiddleware(reproducer),
		middleware.NewDeadLetterErrorMiddleware(eventRepo),
	)

	// Configure breaker-based pause duration and middleware when breaker is enabled
	if appConfig.Breaker != nil {
		conf := *appConfig.Breaker
		conf.Name = appConfig.Name + "_" + contextName + "_kafka_consumer_circuit_breaker"
		router.Use(middleware.NewCircuitBreakMiddleware(&conf))
		consumer.SetBreakerPause(conf.Timeout)
	}
	return consumer, nil
}

func provideAuthEventConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	authUserCreated authApp.UserCreatedUseCase,
) (AuthKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "authorization",
		[]event.Topic{
			authDomain.TopicUserCreated,
		},
		func(r *kafkaInfra.Router) {
			r.EventHandle(authDomain.TopicUserCreated, authEvent.NewUserCreatedEventHandler(authUserCreated))
		},
	)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func provideProfileEventConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	profileUserCreated profileApp.UserCreatedUseCase,
	profileRoomCreated profileApp.RoomCreatedUseCase,
	profileRoomshipCreated profileApp.RoomshipCreatedUseCase,
) (ProfileKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "profile",
		[]event.Topic{
			profileDomain.TopicUserCreated,
			profileDomain.TopicRoomCreated,
			profileDomain.TopicRoomshipCreated,
		},
		func(r *kafkaInfra.Router) {
			r.EventHandle(profileDomain.TopicUserCreated, profileEvent.NewUserCreatedEventHandler(profileUserCreated))
			r.EventHandle(profileDomain.TopicRoomCreated, profileEvent.NewRoomCreatedEventHandler(profileRoomCreated))
			r.EventHandle(profileDomain.TopicRoomshipCreated, profileEvent.NewRoomshipCreatedEventHandler(profileRoomshipCreated))
		},
	)
	return consumer, err
}

func provideChatEventConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomshipCreated chatApp.RoomshipCreatedUseCase,
	chatFriendshipCreated chatApp.FriendshipCreatedUseCase,
	chatUndeliveredMessagesPushRequested chatApp.UndeliveredMessagesPushRequestedUseCase,
) (ChatKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "chat",
		[]event.Topic{
			chatDomain.TopicUserCreated,
			chatDomain.TopicRoomCreated,
			chatDomain.TopicRoomshipCreated,
			chatDomain.TopicFriendshipCreated,
			chatDomain.TopicUndeliveredMessagesPushRequested,
		},
		func(r *kafkaInfra.Router) {
			r.EventHandle(chatDomain.TopicUserCreated, chatEvent.NewUserCreatedEventHandler(chatUserCreated))
			r.EventHandle(chatDomain.TopicRoomCreated, chatEvent.NewRoomCreatedEventHandler(chatRoomCreated))
			r.EventHandle(chatDomain.TopicRoomshipCreated, chatEvent.NewRoomshipCreatedEventHandler(chatRoomshipCreated))
			r.EventHandle(chatDomain.TopicFriendshipCreated, chatEvent.NewFriendshipCreatedEventHandler(chatFriendshipCreated))
			r.EventHandle(chatDomain.TopicUndeliveredMessagesPushRequested, chatEvent.NewUndeliveredMessagesPushRequestedEventHandler(chatUndeliveredMessagesPushRequested))
		},
	)
	return consumer, err
}

func provideNotificationEventConsumer(
	appConfig *config.App,
	emailAvailable emailServiceAvailable,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	notificationWelcomeEmailNotificationRequested notificationApp.WelcomeEmailNotificationRequestedUseCase,
	notificationSystemMessageNotificationRequested notificationApp.SystemMessageNotificationRequestedUseCase,
	notificationUndeliveredMessagesRequested notificationApp.UndeliveredMessagesNotificationRequestedUseCase,
) (NotificationKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "notification",
		[]event.Topic{
			notificationDomain.TopicWelcomeEmailNotificationRequested,
			notificationDomain.TopicSystemMessageNotificationRequested,
			notificationDomain.TopicUndeliveredMessagesNotificationRequested,
		},
		func(r *kafkaInfra.Router) {
			if emailAvailable {
				r.EventHandle(notificationDomain.TopicWelcomeEmailNotificationRequested, notificationEvent.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested))
			} else {
				zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
			}
			r.EventHandle(notificationDomain.TopicSystemMessageNotificationRequested, notificationEvent.NewSystemMessageNotificationRequestedEventHandler(notificationSystemMessageNotificationRequested))
			r.EventHandle(notificationDomain.TopicUndeliveredMessagesNotificationRequested, notificationEvent.NewUndeliveredMessagesNotificationRequestedEventHandler(notificationUndeliveredMessagesRequested))
		},
	)
	return consumer, err
}

func provideRoomshipEventConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	roomshipUserCreated roomshipApp.UserCreatedUseCase,
	roomshipRoomCreated roomshipApp.RoomCreatedUseCase,
	roomshipMemberRequestAgreed roomshipApp.MemberRequestAgreedUseCase,
	roomshipMemberRequestCreated roomshipApp.MemberRequestCreatedUseCase,
	roomshipRoomshipCreated roomshipApp.RoomshipCreatedUseCase,
) (RoomshipKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "roomship",
		[]event.Topic{
			roomshipDomain.TopicUserCreated,
			roomshipDomain.TopicRoomCreated,
			roomshipDomain.TopicMemberRequestCreated,
			roomshipDomain.TopicMemberRequestAgreed,
			roomshipDomain.TopicRoomshipCreated,
		},
		func(r *kafkaInfra.Router) {
			r.EventHandle(roomshipDomain.TopicUserCreated, roomshipEvent.NewUserCreatedEventHandler(roomshipUserCreated))
			r.EventHandle(roomshipDomain.TopicRoomCreated, roomshipEvent.NewRoomCreatedEventHandler(roomshipRoomCreated))
			r.EventHandle(roomshipDomain.TopicMemberRequestCreated, roomshipEvent.NewMemberRequestCreatedEventHandler(roomshipMemberRequestCreated))
			r.EventHandle(roomshipDomain.TopicMemberRequestAgreed, roomshipEvent.NewMemberRequestAgreedEventHandler(roomshipMemberRequestAgreed))
			r.EventHandle(roomshipDomain.TopicRoomshipCreated, roomshipEvent.NewRoomshipCreatedEventHandler(roomshipRoomshipCreated))
		},
	)
	return consumer, err
}

func provideFriendshipEventConsumer(
	appConfig *config.App,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	friendshipUserCreated friendshipApp.UserCreatedUseCase,
	friendshipFriendRequestAgreed friendshipApp.FriendRequestAgreedUseCase,
	friendshipFriendRequestCreated friendshipApp.FriendRequestCreatedUseCase,
	friendshipFriendshipCreated friendshipApp.FriendshipCreatedUseCase,
) (FriendshipKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, reproducer, eventRepo, "friendship",
		[]event.Topic{
			friendshipDomain.TopicUserCreated,
			friendshipDomain.TopicFriendRequestAgreed,
			friendshipDomain.TopicFriendRequestCreated,
			friendshipDomain.TopicFriendshipCreated,
		},
		func(r *kafkaInfra.Router) {
			r.EventHandle(friendshipDomain.TopicUserCreated, friendshipEvent.NewUserCreatedEventHandler(friendshipUserCreated))
			r.EventHandle(friendshipDomain.TopicFriendRequestAgreed, friendshipEvent.NewFriendRequestAgreedEventHandler(friendshipFriendRequestAgreed))
			r.EventHandle(friendshipDomain.TopicFriendRequestCreated, friendshipEvent.NewFriendRequestCreatedEventHandler(friendshipFriendRequestCreated))
			r.EventHandle(friendshipDomain.TopicFriendshipCreated, friendshipEvent.NewFriendshipCreatedEventHandler(friendshipFriendshipCreated))
		},
	)
	return consumer, err
}
