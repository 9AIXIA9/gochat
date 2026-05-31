//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import (
	"context"
	"encoding/json"
	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"
)

// DeliveryService 抽象了通知的投递能力，通知编排层通过它把决策下发到具体渠道。
type DeliveryService interface {
	Deliver(ctx context.Context, recipient kernel.UserID, action command.Action, rawPayload json.RawMessage) error
}
