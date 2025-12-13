package gin

import (
	"gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
)

func Response(ginContext *gin.Context, code http.BusinessCode) {
	ginContext.JSON(code.ToHTTPCode(), http.NewApiResponse(code))
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
