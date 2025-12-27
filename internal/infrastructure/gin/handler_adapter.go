package gin

import (
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdaptUseCaseToHandler[
	Request any,
	RequestPointer RequestPointers[Request],
	Input kernel.Validatable,
	Output any,
](
	useCase kernel.UseCase[Input, Output],
	validator Validator,
	convertRequestToInput func(RequestPointer) Input,
	handleOutput func(*gin.Context, Output),
) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var request RequestPointer = new(Request)
		if err := request.Bind(ginContext); err != nil {
			zap.L().Error("bind request failed", zap.Error(err))
			Response(ginContext, api.CodeInvalidParam)
			return
		}

		message, err := validator.Validate(ginContext.Request.Context(), request)
		if err != nil {
			Response(ginContext, api.CodeServerError)
			return
		}

		if len(message) != 0 {
			ResponseWithMessage(ginContext, api.CodeInvalidParam, message)
			return
		}

		input := convertRequestToInput(request)
		if err := input.Validate(); err != nil {
			Response(ginContext, api.CodeInvalidParam)
			return
		}

		output, err := useCase.Execute(ginContext.Request.Context(), input)
		if err != nil {
			if myErrors.IsBusinessError(err) {
				ResponseWithMessage(ginContext, api.CodeBusinessError, err.Error())
			} else {
				zap.L().Error(
					"http handler failed",
					zap.String("path", ginContext.FullPath()),
					zap.Error(err),
				)
				Response(ginContext, api.CodeServerError)
			}
			return
		}

		if handleOutput != nil {
			handleOutput(ginContext, output)
		}
	}
}
