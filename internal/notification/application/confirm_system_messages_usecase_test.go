package application_test

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeSystemMessagesStatesUpdaterByMessageIDs struct {
	call func(ctx context.Context, userID kernel.UserID, ids []kernel.MessageID, state domain.MessageState) error
}

func (f fakeSystemMessagesStatesUpdaterByMessageIDs) UpdatesByMessageIDs(ctx context.Context, userID kernel.UserID, ids []kernel.MessageID, state domain.MessageState) error {
	return f.call(ctx, userID, ids, state)
}

func TestConfirmSystemMessagesInput_Validate(t *testing.T) {
	input := &application.ConfirmSystemMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: []kernel.MessageID{fixedMessageID},
	}
	err := input.Validate()
	require.NoError(t, err)

	inputWithEmpty := &application.ConfirmSystemMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: nil,
	}
	err = inputWithEmpty.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewConfirmSystemMessagesUseCase(t *testing.T) {
	useCase, err := application.NewConfirmSystemMessagesUseCase(fakeSystemMessagesStatesUpdaterByMessageIDs{call: func(context.Context, kernel.UserID, []kernel.MessageID, domain.MessageState) error {
		return nil
	}})
	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewConfirmSystemMessagesUseCase(nil)
	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestConfirmSystemMessagesUseCase_Execute(t *testing.T) {
	called := false
	useCase, err := application.NewConfirmSystemMessagesUseCase(fakeSystemMessagesStatesUpdaterByMessageIDs{call: func(ctx context.Context, userID kernel.UserID, ids []kernel.MessageID, state domain.MessageState) error {
		called = true
		require.Nil(t, ctx)
		require.Equal(t, fixedUserID, userID)
		require.Equal(t, []kernel.MessageID{fixedMessageID}, ids)
		require.Equal(t, domain.MessageStateDelivered, state)
		return nil
	}})
	require.NoError(t, err)

	_, err = useCase.Execute(nil, &application.ConfirmSystemMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: []kernel.MessageID{fixedMessageID},
	})
	require.NoError(t, err)
	require.True(t, called)
}
