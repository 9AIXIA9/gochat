//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package event

import "context"

type AsyncPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type IDGenerator interface {
	Generate() ID
}
