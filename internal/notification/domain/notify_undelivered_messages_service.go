package domain

import "gochat/internal/shared/kernel"

func NotifyUndeliveredMessages(
	userID kernel.UserID,
	privateMessagesUndelivered []*PrivateMessage,
	roomMessagesUndelivered []*RoomMessage,
	privateMessageNotifier PrivateMessageNotifier,
	roomMessageNotifier RoomMessageNotifier,
) error {

	for _, message := range privateMessagesUndelivered {
		if err := privateMessageNotifier.NotifyPrivateMessage(message); err != nil {
			return err
		}
		message.state = MessageStateDelivered
	}

	for _, message := range roomMessagesUndelivered {
		if ids, err := roomMessageNotifier.NotifyRoomMessage(message, []kernel.UserID{userID}); err != nil || len(ids) == 0 {
			return err // ids == 0 -> nil
		}
		message.states[userID] = MessageStateDelivered
	}
	return nil
}
