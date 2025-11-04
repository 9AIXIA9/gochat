package kafka

import (
	authorizationDomain "gochat/internal/authorization/domain"
	authorizationKafka "gochat/internal/authorization/port/kafka"
	chatDomain "gochat/internal/chat/domain"
	chatKafka "gochat/internal/chat/port/kafka"
	kafkautil "gochat/internal/infrastructure/kafka"
	notificationUsecase "gochat/internal/notification/application/usecase"
	notificationDomain "gochat/internal/notification/domain"
	notificationKafka "gochat/internal/notification/port/kafka"
	"gochat/internal/shared/event"
)

func NewSubscriber(
	commonConfig *kafkautil.CommonConfig,
	consumerConfig *kafkautil.ConsumerConfig,
	eventIDGenerator event.IDGenerator,
	publisher event.Publisher,
	sendEmailUseCase notificationUsecase.SendEmailUseCase,
	sendMessageUseCase notificationUsecase.SendMessageUseCase,
) (*kafkautil.EventSubscriber, error) {
	kafkaSubscriber, err := kafkautil.NewEventSubscriber(commonConfig, consumerConfig)
	if err != nil {
		return nil, err
	}

	kafkaSubscriber.Subscribe(authorizationDomain.TopicUserCreated, authorizationKafka.NewUserCreatedHandler(eventIDGenerator, publisher))

	kafkaSubscriber.Subscribe(chatDomain.TopicMessageSent, chatKafka.NewMessageSentHandler(eventIDGenerator, publisher))

	kafkaSubscriber.Subscribe(notificationDomain.TopicEmailNotificationRequested, notificationKafka.NewEmailNotificationRequestedHandler(sendEmailUseCase))
	kafkaSubscriber.Subscribe(notificationDomain.TopicMessageNotificationRequested, notificationKafka.NewMessageNotificationRequestedHandler(sendMessageUseCase))

	return kafkaSubscriber, nil
}
