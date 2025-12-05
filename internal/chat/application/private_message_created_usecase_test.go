package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPrivateMessageCreatedInput_Validate(t *testing.T) {
	input := &application.PrivateMessageCreatedInput{
		MessageID: fixedMessageID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyMessageID := &application.PrivateMessageCreatedInput{
		MessageID: "",
	}

	err = inputWithEmptyMessageID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewPrivateMessageCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventCreator(ctrl)
	mockPrivateMessageFinder := mocks.NewMockPrivateMessageFinder(ctrl)

	useCase, err := application.NewPrivateMessageCreatedUseCase(
		mockEventIDGenerator,
		mockPrivateMessageFinder,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewPrivateMessageCreatedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestPrivateMessageCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventCreator(ctrl)
	mockPrivateMessageFinder := mocks.NewMockPrivateMessageFinder(ctrl)

	useCase, err := application.NewPrivateMessageCreatedUseCase(
		mockEventIDGenerator,
		mockPrivateMessageFinder,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockMessage := domain.LoadPrivateMessage(
		fixedMessageID,
		fixedUserID,
		fixedFriendID,
		fixedContent,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockPrivateMessageFinder.EXPECT().FindPrivateMessage(nil, fixedMessageID).Return(mockMessage, nil).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockCreator.EXPECT().CreateUnpublishedEvent(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.PrivateMessageCreatedInput{
		MessageID: fixedMessageID,
	})
	require.NoError(t, err)

	// 发送给自己的消息不发送通知
	mockMessageToSelf := domain.LoadPrivateMessage(
		fixedMessageID,
		fixedUserID,
		fixedUserID,
		fixedContent,
		time.Now().UTC(),
	)

	mockPrivateMessageFinder.EXPECT().FindPrivateMessage(nil, fixedMessageID).Return(mockMessageToSelf, nil).Times(1)

	_, err = useCase.Execute(nil, &application.PrivateMessageCreatedInput{
		MessageID: fixedMessageID,
	})
	require.NoError(t, err)
}
