package httptest

import (
	ginutils "gochat/internal/infrastructure/gin"
	"testing"

	"github.com/gin-gonic/gin"
)

func SetupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	ginutils.SetGlobalEnv(gin.TestMode)
	r := gin.Default()
	return r
}
