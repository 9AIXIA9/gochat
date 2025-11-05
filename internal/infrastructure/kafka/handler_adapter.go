package kafka

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func AdaptUseCaseToHandler[
	Input kernel.Validatable,
	Output any,
](
	useCase kernel.UseCase[Input, Output],
	convertEventToInput func(event.Event) (Input, error),
	handleOutput func(context.Context, Output) error,
	handleError func(context.Context, error) error,
) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		input, err := convertEventToInput(e)
		if err != nil {
			return err
		}

		if err := input.Validate(); err != nil {
			return handleError(ctx, err)
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			return handleError(ctx, err)
		}

		return handleOutput(ctx, output)
	}
}
