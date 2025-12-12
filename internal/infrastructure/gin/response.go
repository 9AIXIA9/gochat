package gin

import (
	"gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
)

//TODO 多添加几种形式的响应封装 防止过长以及粗心导致忘记回复

func Response(ginContext *gin.Context, response *http.ApiResponse) {
	ginContext.JSON(response.Code.ToHTTPCode(), response)
}
