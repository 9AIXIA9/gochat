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
)

var KafkaSet = wire.NewSet(
	provideKafkaProducer,
	provideKafkaConsumers,
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

func consumerGroupID(appConfig *config.App, contextName string) string {
	if appConfig == nil || appConfig.Kafka == nil {
		return ""
	}
	return appConfig.Kafka.GroupID + "_" + contextName
}

func buildKafkaConsumer(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	contextName string,
	topicName string,
	register func(r *kafkaInfra.Router),
) (*kafkaInfra.Consumer, error) {
	if appConfig == nil || appConfig.Kafka == nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: app config kafka is nil")
	}

	consumerConfig := appConfig.Kafka.WithGroupIDSuffix(contextName)
	groupID := consumerGroupID(appConfig, contextName)
	consumerName := fmt.Sprintf("%s.%s.%s.kafka_consumer", appConfig.Name, contextName, topicName)
	breakerName := fmt.Sprintf("%s_%s_%s_kafka_consumer_circuit_breaker", appConfig.Name, contextName, topicName)

	router := kafkaInfra.NewRouter()
	var inboxStore middleware.InboxStore
	if redisClient != nil {
		inboxStore = redisInfra.NewInboxStore(redisClient)
	}
	// Order matters: timeout wraps recover so panic inside timeout goroutine is still recoverable.
	middlewares := []kafkaInfra.Middleware{
		middleware.NewRateLimitMiddleware((*limiter.Limiter)(kafkaLimiter), groupID),
		middleware.NewTimeoutMiddleware(appConfig.Timeout),
		middleware.NewRecoverMiddleware(),
		// inbox middleware ensures idempotence by reserving event ids via the inbox store abstraction
		middleware.NewInboxMiddleware(inboxStore, groupID),
		middleware.NewLoggerMiddleware(),
	}
	if appConfig.OTEL != nil && appConfig.OTEL.Enabled {
		middlewares = append([]kafkaInfra.Middleware{middleware.NewTraceMiddleware(consumerName)}, middlewares...)
	}
	router.Use(middlewares...)

	if appConfig.Breaker != nil {
		conf := *appConfig.Breaker
		conf.Name = breakerName
		router.Use(middleware.NewCircuitBreakMiddleware(&conf))
	}
	register(router)

	consumer, err := kafkaInfra.NewConsumer(consumerConfig, router)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer.SetErrorHandler(
		handler.NewLoggerErrorHandler(),
		middleware.NewRetryErrorMiddleware(reproducer),
		middleware.NewDeadLetterErrorMiddlewareWithNamespace(eventRepo, groupID),
	)

	// Configure breaker-based pause duration and middleware when breaker is enabled
	if appConfig.Breaker != nil {
		conf := *appConfig.Breaker
		conf.Name = breakerName
		consumer.SetBreakerPause(conf.Timeout)
	}
	return consumer, nil
}

func provideKafkaConsumers(
	appConfig *config.App,
	kafkaLimiter *KafkaLimiter,
	redisClient *goredis.Client,
	reproducer *ckafka.Producer,
	eventRepo event.Repository,
	profileUserCreated profileApp.UserCreatedUseCase,
	profileRoomCreated profileApp.RoomCreatedUseCase,
	profileRoomshipCreated profileApp.RoomshipCreatedUseCase,
	chatUserCreated chatApp.UserCreatedUseCase,
	chatRoomCreated chatApp.RoomCreatedUseCase,
	chatRoomshipCreated chatApp.RoomshipCreatedUseCase,
	chatFriendshipCreated chatApp.FriendshipCreatedUseCase,
	roomshipUserCreated roomshipApp.UserCreatedUseCase,
	roomshipRoomCreated roomshipApp.RoomCreatedUseCase,
	roomshipMemberRequestAgreed roomshipApp.MemberRequestAgreedUseCase,
	friendshipUserCreated friendshipApp.UserCreatedUseCase,
	friendshipFriendRequestAgreed friendshipApp.FriendRequestAgreedUseCase,
) (consumers []*kafkaInfra.Consumer, err error) {
	defer func() {
		if err == nil {
			return
		}
		for _, consumer := range consumers {
			if consumer != nil {
				consumer.Close()
			}
		}
	}()

	addConsumer := func(consumer *kafkaInfra.Consumer, buildErr error) error {
		if buildErr != nil {
			return buildErr
		}
		consumers = append(consumers, consumer)
		return nil
	}

	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "profile", string(profileDomain.TopicUserCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(profileDomain.TopicUserCreated, profileEvent.NewUserCreatedEventHandler(profileUserCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "profile", string(profileDomain.TopicRoomCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(profileDomain.TopicRoomCreated, profileEvent.NewRoomCreatedEventHandler(profileRoomCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "profile", string(profileDomain.TopicRoomshipCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(profileDomain.TopicRoomshipCreated, profileEvent.NewRoomshipCreatedEventHandler(profileRoomshipCreated))
	})); err != nil {
		return nil, err
	}

	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "chat", string(chatDomain.TopicUserCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(chatDomain.TopicUserCreated, chatEvent.NewUserCreatedEventHandler(chatUserCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "chat", string(chatDomain.TopicRoomCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(chatDomain.TopicRoomCreated, chatEvent.NewRoomCreatedEventHandler(chatRoomCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "chat", string(chatDomain.TopicRoomshipCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(chatDomain.TopicRoomshipCreated, chatEvent.NewRoomshipCreatedEventHandler(chatRoomshipCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "chat", string(chatDomain.TopicFriendshipCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(chatDomain.TopicFriendshipCreated, chatEvent.NewFriendshipCreatedEventHandler(chatFriendshipCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "roomship", string(roomshipDomain.TopicUserCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(roomshipDomain.TopicUserCreated, roomshipEvent.NewUserCreatedEventHandler(roomshipUserCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "roomship", string(roomshipDomain.TopicRoomCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(roomshipDomain.TopicRoomCreated, roomshipEvent.NewRoomCreatedEventHandler(roomshipRoomCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "roomship", string(roomshipDomain.TopicMemberRequestAgreed), func(r *kafkaInfra.Router) {
		r.EventHandle(roomshipDomain.TopicMemberRequestAgreed, roomshipEvent.NewMemberRequestAgreedEventHandler(roomshipMemberRequestAgreed))
	})); err != nil {
		return nil, err
	}

	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "friendship", string(friendshipDomain.TopicUserCreated), func(r *kafkaInfra.Router) {
		r.EventHandle(friendshipDomain.TopicUserCreated, friendshipEvent.NewUserCreatedEventHandler(friendshipUserCreated))
	})); err != nil {
		return nil, err
	}
	if err = addConsumer(buildKafkaConsumer(appConfig, kafkaLimiter, redisClient, reproducer, eventRepo, "friendship", string(friendshipDomain.TopicFriendRequestAgreed), func(r *kafkaInfra.Router) {
		r.EventHandle(friendshipDomain.TopicFriendRequestAgreed, friendshipEvent.NewFriendRequestAgreedEventHandler(friendshipFriendRequestAgreed))
	})); err != nil {
		return nil, err
	}

	return consumers, nil
}
