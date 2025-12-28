package gin

import (
	"gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Response(ginContext *gin.Context, code api.Code) {
	switch code {
	case api.CodeSuccess:
		ginContext.JSON(api.CodeSuccess.ToHTTPCode(), api.ResponseSuccess)
	case api.CodeServerError:
		ginContext.JSON(api.CodeServerError.ToHTTPCode(), api.ResponseServerError)
	case api.CodeTimeout:
		ginContext.JSON(api.CodeTimeout.ToHTTPCode(), api.ResponseTimeout)
	case api.CodeInvalidParam:
		ginContext.JSON(api.CodeInvalidParam.ToHTTPCode(), api.ResponseInvalidParam)
	case api.CodeUnauthorized:
		ginContext.JSON(api.CodeUnauthorized.ToHTTPCode(), api.ResponseInvalidToken)
	case api.CodeServiceUnavailable:
		ginContext.JSON(api.CodeServiceUnavailable.ToHTTPCode(), api.ResponseServiceUnavailable)
	case api.CodeNotFound:
		ginContext.JSON(api.CodeNotFound.ToHTTPCode(), api.ResponseNotFound)
	case api.CodeBusinessError:
		ginContext.JSON(api.CodeBusinessError.ToHTTPCode(), api.ResponseBusinessError)
	default:
		zap.L().Debug("Unhandled business code", zap.Int("code", int(code)))
		ginContext.JSON(code.ToHTTPCode(), api.NewResponse(code))
	}
}

func ResponseSuccess(ginContext *gin.Context) {
	ginContext.JSON(
		api.CodeSuccess.ToHTTPCode(),
		api.ResponseSuccess,
	)
}

func ResponseSuccessWithData(ginContext *gin.Context, data any) {
	ginContext.JSON(
		api.CodeSuccess.ToHTTPCode(),
		api.NewResponseWithData(data),
	)
}

func ResponseWithMessage(ginContext *gin.Context, code api.Code, message string) {
	ginContext.JSON(
		code.ToHTTPCode(),
		api.NewResponseWithMessage(code, message),
	)
}
