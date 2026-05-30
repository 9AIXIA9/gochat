package application_test

import (
	"context"
	"gochat/internal/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTestEvent(id string) event.Event {
	return event.LoadStandardEvent(
		event.ID(id),
		"aggregate-1",
		time.Unix(0, 0).UTC(),
		"test.topic",
		[]byte("payload"),
		nil,
	)
}

func TestNewUnpublishedEventsCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockPublisher := eventMock.NewMockAsyncPublisher(ctrl)
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
	
	mockPublisher := eventMock.NewMockAsyncPublisher(ctrl)
	mockLister := eventMock.NewMockUnpublishedEventsLister(ctrl)
	
	useCase, err := application.NewUnpublishedEventsCreatedUseCase(
		mockPublisher,
		mockLister,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)
	
	firstBatch := []event.Event{newTestEvent("event-1"), newTestEvent("event-2")}
	secondBatch := []event.Event{newTestEvent("event-3")}
	
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(gomock.Any(), time.Minute, 2).Return(firstBatch, nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), firstBatch[0]).Return(nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), firstBatch[1]).Return(nil),
		mockLister.EXPECT().ListUnpublishedEvents(gomock.Any(), time.Minute, 2).Return(secondBatch, nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), secondBatch[0]).Return(nil),
	)
	
	useCase, err = application.NewUnpublishedEventsCreatedUseCaseWithOptions(
		mockPublisher,
		mockLister,
		time.Minute,
		2,
		time.Minute,
	)
	require.NoError(t, err)
	
	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)
	
	//无事件情况
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(gomock.Any(), time.Minute, 2).Return([]event.Event{}, nil),
	)
	
	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)
	
	// 超时情况
	gomock.InOrder(
		mockLister.EXPECT().ListUnpublishedEvents(gomock.Any(), time.Minute, 2).Return([]event.Event{newTestEvent("event-4")}, nil),
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded),
	)
	_, err = useCase.Execute(context.Background(), nil)
	require.NoError(t, err)
}
