package event

import (
	"context"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToEventHandler[
	SpecialEvent SpecificEvent,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertEvToSpecialEv func(Event) (SpecialEvent, error),
	convertEvToInput func(SpecialEvent) Input,
	handleOutput func(context.Context, Output),
	handleError func(context.Context, error),
) Handler {
	return HandlerFunc(func(ctx context.Context, e Event) (err error) {
		defer func() {
			if err != nil && handleError != nil {
				handleError(ctx, err)
			}
		}()
		specialEvent, err := convertEvToSpecialEv(e)
		if err != nil {
			return err
		}

		input := convertEvToInput(specialEvent)

		if err := input.Validate(); err != nil {
			return err
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return err
		}

		if handleOutput == nil {
			return nil
		}
		handleOutput(ctx, output)
		return nil
	})
}
