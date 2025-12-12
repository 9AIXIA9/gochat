package event_test

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	timeTolerance                = 50 * time.Millisecond
	fixedEventID     event.ID    = "event-1234"
	fixedAggregateID kernel.ID   = "aggregate-5678"
	fixedTopic       event.Topic = "test.topic"
	fixedPayload     string      = "This is a test payload."
)

func TestLoadStandardEvent(t *testing.T) {
	ev := event.LoadStandardEvent(
		fixedEventID,
		fixedAggregateID,
		time.Now().UTC(),
		fixedTopic,
		[]byte(fixedPayload),
	)
	require.NotNil(t, ev)
	assert.WithinDuration(t, time.Now().UTC(), ev.OccurredAt(), timeTolerance)
	assert.Equal(t, fixedEventID, ev.ID())
	assert.Equal(t, fixedAggregateID, ev.AggregateID())
	assert.Equal(t, fixedTopic, ev.Topic())
	assert.Equal(t, []byte(fixedPayload), ev.Payload())
}

func TestLoadStandardEventFromEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	fixedOccurredAt := time.Now().UTC()

	mockEv := mocks.NewMockEvent(ctrl)
	mockEv.EXPECT().ID().Return(fixedEventID).Times(1)
	mockEv.EXPECT().AggregateID().Return(fixedAggregateID).Times(1)
	mockEv.EXPECT().OccurredAt().Return(fixedOccurredAt).Times(1)
	mockEv.EXPECT().Topic().Return(fixedTopic).Times(1)
	mockEv.EXPECT().Payload().Return([]byte(fixedPayload)).Times(1)

	standardEv := event.LoadStandardEventFromEvent(mockEv)
	require.NotNil(t, standardEv)
	assert.Equal(t, fixedEventID, standardEv.ID())
	assert.Equal(t, fixedAggregateID, standardEv.AggregateID())
	assert.Equal(t, fixedOccurredAt, standardEv.OccurredAt())
	assert.Equal(t, fixedTopic, standardEv.Topic())
	assert.Equal(t, []byte(fixedPayload), standardEv.Payload())

	// 已经是 standard event
	newStandardEv := event.LoadStandardEventFromEvent(standardEv)
	require.NotNil(t, newStandardEv)
	assert.Equal(t, fixedEventID, newStandardEv.ID())
	assert.Equal(t, fixedAggregateID, newStandardEv.AggregateID())
	assert.Equal(t, fixedOccurredAt, newStandardEv.OccurredAt())
	assert.Equal(t, fixedTopic, newStandardEv.Topic())
	assert.Equal(t, []byte(fixedPayload), newStandardEv.Payload())
}

func TestNewStandardEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	mockIDGen.EXPECT().Generate().Return(fixedEventID).Times(1)

	ev := event.NewStandardEvent(
		fixedAggregateID,
		fixedTopic,
		[]byte(fixedPayload),
		mockIDGen,
	)

	require.NotNil(t, ev)
	assert.Equal(t, fixedEventID, ev.ID())
	assert.Equal(t, fixedAggregateID, ev.AggregateID())
	assert.WithinDuration(t, time.Now().UTC(), ev.OccurredAt(), timeTolerance)
	assert.Equal(t, fixedTopic, ev.Topic())
	assert.Equal(t, []byte(fixedPayload), ev.Payload())
}
