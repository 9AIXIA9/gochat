package event

import (
	"context"
)

type Publisher interface {
	PublishEvents(ctx context.Context, events []Event) error
}
