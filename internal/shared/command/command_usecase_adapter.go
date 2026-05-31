package command

import (
	"context"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

func AdaptUsecaseToCommandHandler[
	SpecialCommand SpecificCommand,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	convertCommandToSpecialCommand func(Command) (SpecialCommand, error),
	convertCommandToInput func(SpecialCommand) Input,
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
			return err
		}
		return nil
	})
}
