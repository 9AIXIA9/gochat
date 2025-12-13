package gin

import (
	"gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Response(ginContext *gin.Context, code http.BusinessCode) {
	switch code {
	case http.CodeSuccess:
		ginContext.JSON(http.CodeSuccess.ToHTTPCode(), http.ResponseSuccess)
	case http.CodeServerError:
		ginContext.JSON(http.CodeServerError.ToHTTPCode(), http.ResponseServerError)
	case http.CodeTimeout:
		ginContext.JSON(http.CodeTimeout.ToHTTPCode(), http.ResponseTimeout)
	case http.CodeInvalidParam:
		ginContext.JSON(http.CodeInvalidParam.ToHTTPCode(), http.ResponseInvalidParam)
	case http.CodeInvalidToken:
		ginContext.JSON(http.CodeInvalidToken.ToHTTPCode(), http.ResponseInvalidToken)
	default:
		zap.L().Debug("Unhandled business code", zap.Int("code", int(code)))
		ginContext.JSON(code.ToHTTPCode(), http.NewApiResponse(code))
	}
}

func ResponseSuccess(ginContext *gin.Context) {
	ginContext.JSON(
		http.CodeSuccess.ToHTTPCode(),
		http.ResponseSuccess,
	)
}

func ResponseSuccessWithData(ginContext *gin.Context, data any) {
	ginContext.JSON(
		http.CodeSuccess.ToHTTPCode(),
		http.NewApiResponseWithData(data),
	)
}

func ResponseWithMessage(ginContext *gin.Context, code http.BusinessCode, message string) {
	ginContext.JSON(
		code.ToHTTPCode(),
		http.NewApiResponseWithMessage(code, message),
	)
}
