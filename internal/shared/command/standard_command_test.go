package command_test

import (
	"gochat/internal/shared/command"
	"gochat/internal/shared/command/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	timeTolerance                   = 50 * time.Millisecond
	fixedCommandID   command.ID     = "command-1234"
	fixedAggregateID kernel.ID      = "aggregate-5678"
	fixedAction      command.Action = "test.action"
	fixedPayload     string         = "This is a test payload."
)

func TestLoadStandardCommand(t *testing.T) {
	headers := map[string]string{
		"header1": "value1",
		"header2": "value2",
	}
	com := command.LoadStandardCommand(
		fixedCommandID,
		fixedAggregateID,
		time.Now().UTC(),
		fixedAction,
		[]byte(fixedPayload),
		headers,
	)
	require.NotNil(t, com)
	assert.WithinDuration(t, time.Now().UTC(), com.OccurredAt(), timeTolerance)
	assert.Equal(t, fixedCommandID, com.ID())
	assert.Equal(t, fixedAggregateID, com.AggregateID())
	assert.Equal(t, fixedAction, com.Action())
	assert.Equal(t, []byte(fixedPayload), com.Payload())
	assert.Equal(t, headers, com.Headers())
}

func TestLoadStandardCommandFromCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	fixedOccurredAt := time.Now().UTC()

	mockCommand := mocks.NewMockCommand(ctrl)
	gomock.InOrder(
		mockCommand.EXPECT().ID().Return(fixedCommandID).Times(1),
		mockCommand.EXPECT().AggregateID().Return(fixedAggregateID).Times(1),
		mockCommand.EXPECT().OccurredAt().Return(fixedOccurredAt).Times(1),
		mockCommand.EXPECT().Action().Return(fixedAction).Times(1),
		mockCommand.EXPECT().Payload().Return([]byte(fixedPayload)).Times(1),
		mockCommand.EXPECT().Headers().Return(map[string]string{}).Times(1),
	)

	standardCommand := command.LoadStandardCommandFromCommand(mockCommand)
	require.NotNil(t, standardCommand)
	assert.Equal(t, fixedCommandID, standardCommand.ID())
	assert.Equal(t, fixedAggregateID, standardCommand.AggregateID())
	assert.Equal(t, fixedOccurredAt, standardCommand.OccurredAt())
	assert.Equal(t, fixedAction, standardCommand.Action())
	assert.Equal(t, []byte(fixedPayload), standardCommand.Payload())

	// 已经是 standard command
	newStandardCommand := command.LoadStandardCommandFromCommand(standardCommand)
	require.NotNil(t, newStandardCommand)
	assert.Equal(t, fixedCommandID, newStandardCommand.ID())
	assert.Equal(t, fixedAggregateID, newStandardCommand.AggregateID())
	assert.Equal(t, fixedOccurredAt, newStandardCommand.OccurredAt())
	assert.Equal(t, fixedAction, newStandardCommand.Action())
	assert.Equal(t, []byte(fixedPayload), newStandardCommand.Payload())
}

func TestNewStandardCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	mockIDGen.EXPECT().Generate().Return(fixedCommandID).Times(1)

	com := command.NewStandardCommand(
		fixedAggregateID,
		fixedAction,
		[]byte(fixedPayload),
		mockIDGen,
	)

	require.NotNil(t, com)
	assert.Equal(t, fixedCommandID, com.ID())
	assert.Equal(t, fixedAggregateID, com.AggregateID())
	assert.WithinDuration(t, time.Now().UTC(), com.OccurredAt(), timeTolerance)
	assert.Equal(t, fixedAction, com.Action())
	assert.Equal(t, []byte(fixedPayload), com.Payload())
}

func TestStandardCommand_AddHeader(t *testing.T) {
	com := command.LoadStandardCommand(
		fixedCommandID,
		fixedAggregateID,
		time.Now().UTC(),
		fixedAction,
		[]byte(fixedPayload),
		nil,
	)

	com.AddHeader("key1", "value1")
	assert.Equal(t, "value1", com.Headers()["key1"])

	com.AddHeader("key2", "value2")
	assert.Equal(t, "value2", com.Headers()["key2"])
}

func TestStandardCommand_AddHeaders(t *testing.T) {
	com := command.LoadStandardCommand(
		fixedCommandID,
		fixedAggregateID,
		time.Now().UTC(),
		fixedAction,
		[]byte(fixedPayload),
		nil,
	)

	headersToAdd := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	com.AddHeaders(headersToAdd)

	for k, v := range headersToAdd {
		assert.Equal(t, v, com.Headers()[k])
	}
}
