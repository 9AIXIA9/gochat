package gin

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/shared/http"
)

func Response(ginContext *gin.Context, response *http.ApiResponse) {
	ginContext.JSON(response.Code.ToHTTPCode(), response)
}
