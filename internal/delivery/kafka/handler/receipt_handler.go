package handler

import (
	"context"
	"fmt"

	"gochat/internal/gateway/core"
	"gochat/internal/shared/command"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/kernel"
)

const ActionPushCommandReceipt command.Action = "push_command_receipt"

func NewReceiptGatewayHandler(gateway contract.GatewayService) command.ReceiptHandler {
	return command.ReceiptHandlerFunc(func(ctx context.Context, r command.Receipt) error {
		if gateway == nil {
			return fmt.Errorf("gateway service is nil")
		}

		payload, err := r.Marshal()
		if err != nil {
			return err
		}

		envelop := core.NewEnvelopWithClientMessageID(ActionPushCommandReceipt, payload, kernel.MessageID(r.Headers()["client_message_id"]))

		return gateway.PushToUser(ctx, kernel.UserID(r.AggregateID()), envelop)
	})
}
