package gin

import (
	"context"
	"fmt"
	"gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"runtime"
	"time"

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
	timeoutOptions ...time.Duration,
) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var request RequestPointer = new(Request)
		err := request.Bind(ginContext)
		if err != nil {
			zap.L().Error("bind request failed", zap.Error(err))
			Response(ginContext, http.ResponseInvalidParam)
			return
		}

		message, err := validator.Validate(ginContext.Request.Context(), request)
		if err != nil {
			Response(ginContext, http.ResponseServerError)
			return
		}

		if len(message) != 0 {
			Response(ginContext, http.NewApiResponseWithMessage(http.CodeInvalidParam, message))
			return
		}

		input := convertRequestToInput(request)
		if err := input.Validate(); err != nil {
			Response(ginContext, http.ResponseInvalidParam)
			return
		}

		var ctx context.Context

		if len(timeoutOptions) == 0 {
			ctx = ginContext.Request.Context()
		} else {
			ctxWithTimeout, cancel := context.WithTimeout(ginContext.Request.Context(), timeoutOptions[0])
			ctx = ctxWithTimeout
			defer cancel()
		}

		errChan := make(chan error, 1)
		outputChan := make(chan Output, 1)
		defer func() {
			close(outputChan)
			close(errChan)
		}()

		go func() {
			defer func() {
				if r := recover(); r != nil {
					stack := make([]byte, 4096)
					length := runtime.Stack(stack, false)
					select {
					case errChan <- fmt.Errorf("panic: %v\nstack: %s", r, stack[:length]):
					case <-ctx.Done():
					}
				}
			}()

			output, err := useCase.Execute(ctx, input)
			if err != nil {
				select {
				case errChan <- err:
				case <-ctx.Done():
				}
				return
			}
			outputChan <- output
		}()

		select {
		case <-ctx.Done():
			Response(ginContext, http.ResponseTimeout)
			return
		case err := <-errChan:
			handleError(ginContext, err)
		case output := <-outputChan:
			handleOutput(ginContext, output)
		}
	}
}
