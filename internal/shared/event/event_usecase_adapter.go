package event

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

func AdaptUsecaseToEventHandler[
	SpecialEvent SpecificEvent,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertEvToSpecialEv func(Event) (SpecialEvent, error),
	convertEvToInput func(SpecialEvent) Input,
) Handler {
	return HandlerFunc(func(ctx context.Context, e Event) (err error) {
		specialEvent, err := convertEvToSpecialEv(e)
		if err != nil {
			return err
		}

		input := convertEvToInput(specialEvent)

		if err := input.Validate(); err != nil {
			return err
		}

		_, err = usecase.Execute(ctxutil.WithHeaders(ctx, specialEvent.Headers()), input)
		if err != nil {
			return err
		}
		return nil
	})
}
