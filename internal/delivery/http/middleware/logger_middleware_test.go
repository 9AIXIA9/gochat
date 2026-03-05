package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerMiddleware(t *testing.T) {
	t.Parallel()

	router := httptestutil.NewTestRouter(t)
	router.Use(middlewareHTTP.NewLoggerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
