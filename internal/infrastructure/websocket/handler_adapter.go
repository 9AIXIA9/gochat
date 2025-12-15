package websocket

import (
	"context"
	"gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
)

func AdaptUsecaseToHandler[
	RequestData any,
	Input kernel.Validatable,
	Output any,
](
	usecase kernel.UseCase[Input, Output],
	validator Validator,
	bindRequestData func(context.Context, []byte) (RequestData, error),
	mapInput func(RequestData) Input,
	handleError func(context.Context, error) *api.Response,
) Handler {
	return HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
		reqData, err := bindRequestData(ctx, data)
		if err != nil {
			return api.NewResponse(api.CodeInvalidParam)
		}

		message, err := validator.Validate(ctx, reqData)
		if err != nil {
			return api.NewResponse(api.CodeServerError)
		}

		if len(message) != 0 {
			return api.NewResponseWithMessage(api.CodeInvalidParam, message)
		}

		input := mapInput(reqData)

		if err := input.Validate(); err != nil {
			return api.NewResponse(api.CodeInvalidParam)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			if handleError != nil {
				return handleError(ctx, err)
			}
			return api.NewResponse(api.CodeServerError)
		}
		return api.NewResponseWithData(output)
	})
}
