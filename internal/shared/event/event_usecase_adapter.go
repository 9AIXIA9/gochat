package event

import (
	"context"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToHandler[
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
	return HandlerFunc(func(ctx context.Context, e Event) error {
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
			if handleError == nil {
				return err
			}
			handleError(ctx, err)
			return err
		}

		if handleOutput == nil {
			return nil
		}
		handleOutput(ctx, output)
		return nil
	})
}
