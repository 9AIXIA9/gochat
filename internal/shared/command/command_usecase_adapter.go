package command

import (
	"context"
	"errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"go.uber.org/zap"
)

func AdaptUsecaseToCommandHandler[
	SpecialCommand SpecificCommand,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertCommandToSpecialCommand func(Command) (SpecialCommand, error),
	convertCommandToInput func(SpecialCommand) Input,
	publisher ReceiptAsyncPublisher,
	generator ReceiptIDGenerator,
) Handler {
	return HandlerFunc(func(ctx context.Context, e Command) (err error) {
		specialCommand, err := convertCommandToSpecialCommand(e)
		if err != nil {
			return err
		}

		input := convertCommandToInput(specialCommand)
		if err := input.Validate(); err != nil {
			return err
		}

		_, err = usecase.Execute(ctxutil.WithHeaders(ctx, specialCommand.Headers()), input)
		if err != nil {
			receipt := NewStandardReceiptFromCommand(
				specialCommand,
				StatusFailed,
				generator,
			)
			if publishErr := publisher.Publish(ctx, receipt); publishErr != nil {
				zap.L().Error(
					"Failed to publish command error receipt",
					zap.String("command_id", e.ID().String()),
					zap.String("command", e.Action().String()),
					zap.String("aggregate_id", e.AggregateID().String()),
					zap.Error(publishErr),
				)
				return errors.Join(publishErr, err)
			}
			return err
		}

		receipt := NewStandardReceiptFromCommand(
			specialCommand,
			StatusSucceeded,
			generator,
		)
		if publishErr := publisher.Publish(ctx, receipt); publishErr != nil {
			zap.L().Error(
				"Failed to publish command succeeded receipt",
				zap.String("command_id", e.ID().String()),
				zap.String("command", e.Action().String()),
				zap.String("aggregate_id", e.AggregateID().String()),
				zap.Error(publishErr),
			)
		}
		return nil
	})
}
