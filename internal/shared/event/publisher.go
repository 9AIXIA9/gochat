//go:generate mockgen -source=publisher.go -destination=./mocks/mock_publisher.go -package=mocks
package event

import "context"

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
