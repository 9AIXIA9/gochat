package application_test

import (
	"context"
	"fmt"
	"gochat/internal/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
		mockLister.EXPECT().ListUnpublishedEvents(context.Background(), time.Minute).Return(mockEvents, nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(fixedEventLen).Return(nil),
	)

	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)

	//无事件情况
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(context.Background(), time.Minute).Return([]event.Event{}, nil),
	)

	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)

	// 超时情况
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(context.Background(), time.Minute).Return(mockEvents, nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded),
	)
	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)
}

func TestUnpublishedEventsCreatedUseCase_Execute_ConcurrentReentrySkipped(t *testing.T) {
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

	block := make(chan struct{})
	started := make(chan struct{})
	var startOnce sync.Once

	mockGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	mockEvent := event.NewStandardEvent(
		"aggregate-0",
		"event-topic-0",
		nil,
		mockGenerator,
	)

	mockLister.EXPECT().ListUnpublishedEvents(gomock.Any(), time.Minute).DoAndReturn(
		func(_ context.Context, _ time.Duration) ([]event.Event, error) {
			startOnce.Do(func() { close(started) })
			<-block
			return []event.Event{mockEvent}, nil
		},
	).Times(1)

	var publishCount atomic.Int32
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ event.Event) error {
			publishCount.Add(1)
			return nil
		},
	).Times(1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, callErr := useCase.Execute(context.Background(), nil)
		require.NoError(t, callErr)
	}()

	<-started

	go func() {
		defer wg.Done()
		_, callErr := useCase.Execute(context.Background(), nil)
		require.NoError(t, callErr)
	}()

	close(block)
	wg.Wait()

	require.Equal(t, int32(1), publishCount.Load())
}
