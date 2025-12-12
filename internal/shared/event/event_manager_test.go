package event_test

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/event/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedTestEventTimes = 10

func TestNewEventManager(t *testing.T) {
	manager := event.NewEventManager()
	require.NotNil(t, manager)
	assert.Nil(t, manager.GetEvents())
}

func TestManager_RecordAndGetEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manager := event.NewEventManager()
	require.NotNil(t, manager)
	assert.Nil(t, manager.GetEvents())

	//添加事件
	mockEv := mocks.NewMockSpecificEvent(ctrl)
	for i := 0; i < fixedTestEventTimes; i++ {
		manager.RecordEvent(mockEv)
	}
	evs := manager.GetEvents()
	require.NotNil(t, evs)
	assert.Equal(t, fixedTestEventTimes, len(evs))

	//再次获取，应该为空
	evs = manager.GetEvents()
	assert.Nil(t, evs)
}
