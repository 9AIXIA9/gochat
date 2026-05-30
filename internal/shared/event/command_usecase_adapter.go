package event

import (
	"context"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToCommandHandler[
	SpecialEvent SpecificEvent,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertEvToSpecialEv func(Event) (SpecialEvent, error),
	convertEvToInput func(SpecialEvent) Input,
	handleOutput func(context.Context, SpecialEvent, Output),
	handleError func(context.Context, SpecialEvent, error),
) Handler {
	return HandlerFunc(func(ctx context.Context, e Event) (err error) {
		specialEvent, err := convertEvToSpecialEv(e)
		if err != nil {
			return err
		}

		input := convertEvToInput(specialEvent)

		defer func() {
			if err != nil && handleError != nil {
				handleError(ctx, specialEvent, err)
			}
		}()

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
		handleOutput(ctx, specialEvent, output)
		return nil
	})
}
