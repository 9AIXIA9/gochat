package di

import (
	"fmt"
	"gochat/config"
	chatApp "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	chatEvent "gochat/internal/chat/port/event"
	"gochat/internal/delivery/kafka/handler"
	"gochat/internal/delivery/kafka/middleware"
	friendshipApp "gochat/internal/friendship/application"
	friendshipDomain "gochat/internal/friendship/domain"
	friendshipEvent "gochat/internal/friendship/port/event"
	kafkaInfra "gochat/internal/infrastructure/kafka"
	redisInfra "gochat/internal/infrastructure/redis"
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
	goredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"go.uber.org/zap"
)

var KafkaSet = wire.NewSet(
	provideKafkaProducer,
	provideKafkaConsumers,
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

func provideKafkaConsumers(
	profileConsumer ProfileKafkaConsumer,
	chatConsumer ChatKafkaConsumer,
	notificationConsumer NotificationKafkaConsumer,
	roomshipConsumer RoomshipKafkaConsumer,
	friendshipConsumer FriendshipKafkaConsumer,
) []*kafkaInfra.Consumer {
	return []*kafkaInfra.Consumer{
		profileConsumer,
		chatConsumer,
		notificationConsumer,
		roomshipConsumer,
		friendshipConsumer,
	}
}

// common builder to reduce duplication across contexts
func buildKafkaConsumer(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	contextName string,
	register func(r *kafkaInfra.Router),
) (*kafkaInfra.Consumer, error) {
	router := kafkaInfra.NewRouter()
	inboxStore := redisInfra.NewInboxStore(redisClient)
	// Order matters: timeout wraps recover so panic inside timeout goroutine is still recoverable.
	middlewares := []kafkaInfra.Middleware{
		middleware.NewRateLimitMiddleware((*limiter.Limiter)(kafkaLimiter)),
		middleware.NewTimeoutMiddleware(appConfig.Timeout),
		middleware.NewRecoverMiddleware(),
		// inbox middleware ensures idempotence by reserving event ids via the inbox store abstraction
		middleware.NewInboxMiddleware(inboxStore),
		middleware.NewLoggerMiddleware(),
	}
	if appConfig.OTEL != nil && appConfig.OTEL.Enabled {
		middlewares = append([]kafkaInfra.Middleware{middleware.NewTraceMiddleware(appConfig.Name + "." + contextName + ".kafka_consumer")}, middlewares...)
	}
	router.Use(middlewares...)

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
		consumer.SetBreakerPause(conf.Timeout)
	}
	return consumer, nil
}

func provideProfileEventConsumer(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	profileUserCreated profileApp.UserCreatedUseCase,
	profileRoomCreated profileApp.RoomCreatedUseCase,
	profileRoomshipCreated profileApp.RoomshipCreatedUseCase,
) (ProfileKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "profile",
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
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomshipCreated chatApp.RoomshipCreatedUseCase,
	chatFriendshipCreated chatApp.FriendshipCreatedUseCase,
) (ChatKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "chat",
		func(r *kafkaInfra.Router) {
			r.EventHandle(chatDomain.TopicUserCreated, chatEvent.NewUserCreatedEventHandler(chatUserCreated))
			r.EventHandle(chatDomain.TopicRoomCreated, chatEvent.NewRoomCreatedEventHandler(chatRoomCreated))
			r.EventHandle(chatDomain.TopicRoomshipCreated, chatEvent.NewRoomshipCreatedEventHandler(chatRoomshipCreated))
			r.EventHandle(chatDomain.TopicFriendshipCreated, chatEvent.NewFriendshipCreatedEventHandler(chatFriendshipCreated))
		},
	)
	return consumer, err
}

func provideNotificationEventConsumer(
	appConfig *config.App,
	emailAvailable emailServiceAvailable,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	notificationWelcomeEmailNotificationRequested notificationApp.WelcomeEmailNotificationRequestedUseCase,
	notificationSystemMessageNotificationRequested notificationApp.SystemMessageNotificationRequestedUseCase,
) (NotificationKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "notification",
		func(r *kafkaInfra.Router) {
			if !appConfig.Email.Enable {
				zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email notifier is disabled in config")
			} else if emailAvailable {
				r.EventHandle(notificationDomain.TopicWelcomeEmailNotificationRequested, notificationEvent.NewWelcomeEmailNotificationRequestedEventHandler(notificationWelcomeEmailNotificationRequested))
			} else {
				zap.L().Info("Skipping subscription to WelcomeEmailNotificationRequested topic as email dialer is not connected")
			}
			r.EventHandle(notificationDomain.TopicSystemMessageNotificationRequested, notificationEvent.NewSystemMessageNotificationRequestedEventHandler(notificationSystemMessageNotificationRequested))
		},
	)
	return consumer, err
}

func provideRoomshipEventConsumer(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	roomshipUserCreated roomshipApp.UserCreatedUseCase,
	roomshipRoomCreated roomshipApp.RoomCreatedUseCase,
	roomshipMemberRequestAgreed roomshipApp.MemberRequestAgreedUseCase,
	roomshipMemberRequestCreated roomshipApp.MemberRequestCreatedUseCase,
) (RoomshipKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "roomship",
		func(r *kafkaInfra.Router) {
			r.EventHandle(roomshipDomain.TopicUserCreated, roomshipEvent.NewUserCreatedEventHandler(roomshipUserCreated))
			r.EventHandle(roomshipDomain.TopicRoomCreated, roomshipEvent.NewRoomCreatedEventHandler(roomshipRoomCreated))
			r.EventHandle(roomshipDomain.TopicMemberRequestCreated, roomshipEvent.NewMemberRequestCreatedEventHandler(roomshipMemberRequestCreated))
			r.EventHandle(roomshipDomain.TopicMemberRequestAgreed, roomshipEvent.NewMemberRequestAgreedEventHandler(roomshipMemberRequestAgreed))
		},
	)
	return consumer, err
}

func provideFriendshipEventConsumer(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	friendshipUserCreated friendshipApp.UserCreatedUseCase,
	friendshipFriendRequestAgreed friendshipApp.FriendRequestAgreedUseCase,
) (FriendshipKafkaConsumer, error) {
	consumer, err := buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "friendship",
		func(r *kafkaInfra.Router) {
			r.EventHandle(friendshipDomain.TopicUserCreated, friendshipEvent.NewUserCreatedEventHandler(friendshipUserCreated))
			r.EventHandle(friendshipDomain.TopicFriendRequestAgreed, friendshipEvent.NewFriendRequestAgreedEventHandler(friendshipFriendRequestAgreed))
		},
	)
	return consumer, err
}
