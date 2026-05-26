package application_test

import (
	"fmt"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedRoomMembersCount = 10

func TestSendRoomMessageInput_Validate(t *testing.T) {
	input := &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptySenderID := &application.SendRoomMessageInput{
		SenderID: "",
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	}

	err = inputWithEmptySenderID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRecipientID := &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   "",
		Content:  fixedContent,
	}

	err = inputWithEmptyRecipientID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewSendRoomMessageUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockMessageCreator := mocks.NewMockRoomMessageCreator(ctrl)
	mockNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	finder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewSendRoomMessageUseCase(
		mockMessageIDGenerator,
		mockNotifier,
		finder,
		mockMessageCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSendRoomMessageUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSendRoomMessageUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockMessageCreator := mocks.NewMockRoomMessageCreator(ctrl)
	mockNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	finder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewSendRoomMessageUseCase(
		mockMessageIDGenerator,
		mockNotifier,
		finder,
		mockMessageCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRoomships := make([]*domain.Roomship, fixedRoomMembersCount)
	mockMembers := make([]kernel.UserID, 0, fixedRoomMembersCount)
	for i := 0; i < fixedRoomMembersCount; i++ {
		mockRoomships[i] = domain.LoadRoomship(
			domain.RoomshipID(fmt.Sprintf("roomship-%d", i)),
			kernel.UserID(fmt.Sprintf("user-%d", i)),
			fixedRoomID,
		)
		mockMembers = append(mockMembers, mockRoomships[i].UserID())
	}

	mockRoomships = append(mockRoomships, domain.LoadRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
	))
	mockRoomshipsWithoutSender := mockRoomships[:len(mockRoomships)-1]

	// 统一按真实调用顺序设置期望：
	// 1) finder -> 2) Generate -> 3) Notify(若有其他成员) -> 4) Create
	singleRoomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
	)
	gomock.InOrder(
		// 正常情况
		finder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockRoomships, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
		mockNotifier.EXPECT().Notify(nil, gomock.Any(), gomock.InAnyOrder(mockMembers)).Return(nil).Times(1),

		// 不是成员
		finder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockRoomshipsWithoutSender, nil).Times(1),

		// 房间不存在
		finder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return([]*domain.Roomship{}, nil).Times(1),

		// 房间只有自己
		finder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return([]*domain.Roomship{singleRoomship}, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	// 正常情况
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.NoError(t, err)

	// 不是成员
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.ErrorIs(t, err, domain.ErrNotMember)

	// 房间不存在
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.ErrorIs(t, err, domain.ErrRoomNotFound)

	// 房间只有自己
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.NoError(t, err)

	// 创建消息失败
	gomock.InOrder(
		finder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockRoomships, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(myErrors.ErrDuplicatedKey).Times(1),
	)

	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.ErrorIs(t, err, myErrors.ErrDuplicatedKey)
}
