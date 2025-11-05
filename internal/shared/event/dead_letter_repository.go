package event

import "context"

type DeadLetterSaver interface {
	SaveDeadLetter(ctx context.Context, event Event, reason error) error
}
