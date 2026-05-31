//go:generate mockgen -source=command_publisher.go -destination=./mocks/mock_command_publisher.go -package=mocks
package command

import "context"

type SyncPublisher interface {
	Publish(ctx context.Context, command Command) error
}
