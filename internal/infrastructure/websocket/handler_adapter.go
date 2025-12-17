package websocket

import (
	"context"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"

	"go.uber.org/zap"
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
			if myErrors.IsBusinessError(err) {
				return api.NewResponseWithMessage(api.CodeSuccess, err.Error())
			}
			zap.L().Error(
				"websocket usecase execute failed",
				zap.String("user_id", utils.GetUserID(ctx).String()),
				zap.String("topic", GetTopic(ctx).String()),
				zap.Error(err),
			)
			return api.NewResponse(api.CodeServerError)
		}
		return api.NewResponseWithData(output)
	})
}
