package command

import (
	"context"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToCommandHandler[
	SpecialCommand SpecificCommand,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertCommandToSpecialCommand func(Command) (SpecialCommand, error),
	convertCommandToInput func(SpecialCommand) Input,
	handleOutput func(context.Context, SpecialCommand, Output),
	handleError func(context.Context, SpecialCommand, error),
) Handler {
	return HandlerFunc(func(ctx context.Context, e Command) (err error) {
		specialCommand, err := convertCommandToSpecialCommand(e)
		if err != nil {
			return err
		}

		input := convertCommandToInput(specialCommand)

		defer func() {
			if err != nil && handleError != nil {
				handleError(ctx, specialCommand, err)
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
		handleOutput(ctx, specialCommand, output)
		return nil
	})
}
