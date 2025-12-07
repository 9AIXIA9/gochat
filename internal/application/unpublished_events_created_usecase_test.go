package application_test

import (
	"fmt"
	"gochat/internal/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedEventLen = 10

func TestNewUnpublishedEventsCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := eventMock.NewMockPublisher(ctrl)
	mockLister := eventMock.NewMockUnpublishedEventsLister(ctrl)

	useCase, err := application.NewUnpublishedEventsCreatedUseCase(
		mockPublisher,
		mockLister,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUnpublishedEventsCreatedUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUnpublishedEventsCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := eventMock.NewMockPublisher(ctrl)
	mockLister := eventMock.NewMockUnpublishedEventsLister(ctrl)

	useCase, err := application.NewUnpublishedEventsCreatedUseCase(
		mockPublisher,
		mockLister,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().Generate().Return(fixedEventID).Times(fixedEventLen)

	mockEvents := make([]event.Event, 0, fixedEventLen)
	for i := 0; i < fixedEventLen; i++ {
		mockEvents = append(mockEvents, event.NewStandardEvent(
			kernel.ID(fmt.Sprintf("aggregate-%d", i)),
			event.Topic(fmt.Sprintf("event-topic-%d", i)),
			nil,
			mockGenerator,
		))
	}

	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(nil).Return(mockEvents, nil),
		mockPublisher.EXPECT().Publish(gomock.Any()).Times(fixedEventLen).Return(nil),
	)

	_, err = useCase.Execute(nil, nil)
	require.NoError(t, err)

	//无事件情况
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(nil).Return([]event.Event{}, nil),
	)

	_, err = useCase.Execute(nil, nil)
	require.NoError(t, err)
}
