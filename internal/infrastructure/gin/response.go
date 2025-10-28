package gin

import (
	"gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
)

func Response(ginContext *gin.Context, response *http.ApiResponse) {
	ginContext.JSON(response.Code.ToHTTPCode(), response)
}
