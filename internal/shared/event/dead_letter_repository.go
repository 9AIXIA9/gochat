package event

import "context"

type DeadLetterRepository interface {
	DeadLetterSaver
}

type DeadLetterSaver interface {
	SaveDeadLetter(ctx context.Context, event Event, reason error) error
}
