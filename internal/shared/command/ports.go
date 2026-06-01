//go:generate mockgen -source ports.go -destination=./mocks/mock_ports.go -package=mocks
package command

import (
	"context"
)

type SyncPublisher interface {
	Publish(ctx context.Context, command Command) error
}

type ReceiptAsyncPublisher interface {
	Publish(ctx context.Context, receipt Receipt) error
}

type IDGenerator interface {
	Generate() ID
}

type ReceiptIDGenerator interface {
	Generate() ReceiptID
}
