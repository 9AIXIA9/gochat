package kafka

import (
	authorizationDomain "gochat/internal/authorization/domain"
	authorizationKafka "gochat/internal/authorization/port/kafka"
	kafkautil "gochat/internal/infrastructure/kafka"
	notificationUsecase "gochat/internal/notification/application/usecase"
	notificationDomain "gochat/internal/notification/domain"
	notificationKafka "gochat/internal/notification/port/kafka"
	"gochat/internal/shared/event"
)

func NewSubscriber(
	commonConfig *kafkautil.CommonConfig,
	consumerConfig *kafkautil.ConsumerConfig,
	publisher event.Publisher,
	sendEmailUseCase notificationUsecase.SendEmailUseCase,
) (*kafkautil.EventSubscriber, error) {
	kafkaSubscriber, err := kafkautil.NewEventSubscriber(commonConfig, consumerConfig)
	if err != nil {
		return nil, err
	}

	kafkaSubscriber.Subscribe(authorizationDomain.TopicUserCreated, authorizationKafka.NewUserCreatedHandler(publisher))
	kafkaSubscriber.Subscribe(notificationDomain.TopicEmailSendingRequested, notificationKafka.NewEmailSendingRequestedHandler(sendEmailUseCase))

	return kafkaSubscriber, nil
}
