package gateway

import (
	"context"
	"encoding/json"
	"gochat/internal/gateway/core"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ domain.DeliveryService = (*DeliverService)(nil)

type DeliverService struct {
	gateway contract.GatewayService
}

func NewDeliverService(gateway contract.GatewayService) *DeliverService {
	return &DeliverService{gateway: gateway}
}

func (s *DeliverService) Deliver(ctx context.Context, recipient kernel.UserID, action event.Topic, rawPayload json.RawMessage) error {
	return s.gateway.PushToUser(ctx, recipient, core.NewActiveDownstreamEnvelop(action, rawPayload))
}
