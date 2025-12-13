package websocket

import (
	"context"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToHandler[
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	inputMapper func(context.Context, []byte) (Input, error),
	handleOutput func(context.Context, Output) ([]byte, error),
	handleError func(context.Context, error),
) Handler {
	return HandlerFunc(func(ctx context.Context, data []byte) ([]byte, error) {
		input, err := inputMapper(ctx, data)
		if err != nil {
			return nil, err
		}

		if err := input.Validate(); err != nil {
			return nil, err
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			if handleError == nil {
				return nil, err
			}
			return nil, err
		}

		if handleOutput == nil {
			return nil, nil
		}

		bytes, err := handleOutput(ctx, output)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	})
}
