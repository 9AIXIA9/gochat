//go:generate mockgen -source=event_publisher.go -destination=./mocks/mock_event_publisher.go -package=mocks
package event

import "context"

type AsyncPublisher interface {
	Publish(ctx context.Context, event Event) error
}
