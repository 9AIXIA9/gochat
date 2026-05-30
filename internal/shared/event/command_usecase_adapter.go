package event

import (
	"context"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

func AdaptUsecaseToCommandHandler[
	SpecialEvent SpecificEvent,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	publisher SyncPublisher,
	convertEvToSpecialEv func(Event) (SpecialEvent, error),
	convertEvToInput func(SpecialEvent) Input,
	handleOutput func(context.Context, SpecialEvent, Output) SpecificEvent,
	handleError func(context.Context, SpecialEvent, error) SpecificEvent,
) Handler {
	return HandlerFunc(func(ctx context.Context, e Event) (err error) {
		specialEvent, err := convertEvToSpecialEv(e)
		if err != nil {
			return err
		}

		input := convertEvToInput(specialEvent)

		defer func() {
			if err != nil && handleError != nil {
				if errorEv := handleError(ctx, specialEvent, err); errorEv != nil {
					if pubErr := publisher.Publish(ctx, errorEv); pubErr != nil {
						zap.L().Error("failed to publish error event", zap.Error(pubErr), zap.String("original_error", err.Error()))
						return
					}
				}

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
		if successEvent := handleOutput(ctx, specialEvent, output); successEvent != nil {
			if pubErr := publisher.Publish(ctx, successEvent); pubErr != nil {
				zap.L().Error("failed to publish success event", zap.Error(pubErr))
				return
			}
		}
		return nil
	})
}
