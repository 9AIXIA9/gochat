package gin

import (
	"gochat/internal/shared/http"
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
	handleError func(*gin.Context, error),
) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var request RequestPointer = new(Request)
		err := request.Bind(ginContext)
		if err != nil {
			zap.L().Error("bind request failed", zap.Error(err))
			Response(ginContext, http.CodeInvalidParam)
			return
		}

		message, err := validator.Validate(ginContext.Request.Context(), request)
		if err != nil {
			Response(ginContext, http.CodeServerError)
			return
		}

		if len(message) != 0 {
			ResponseWithMessage(ginContext, http.CodeInvalidParam, message)
			return
		}

		input := convertRequestToInput(request)
		if err := input.Validate(); err != nil {
			Response(ginContext, http.CodeInvalidParam)
			return
		}

		output, err := useCase.Execute(ginContext.Request.Context(), input)
		if err != nil {
			if handleError != nil {
				handleError(ginContext, err)
			} else {
				zap.L().Error("use case execute failed", zap.Error(err))
				Response(ginContext, http.CodeServerError)
			}
			return
		}

		if handleOutput != nil {
			handleOutput(ginContext, output)
		}
	}
}
