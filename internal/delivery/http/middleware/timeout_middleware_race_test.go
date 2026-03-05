//go:build race

package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewTimeoutMiddleware_RaceBuild(t *testing.T) {
	t.Parallel()

	router := httptestutil.NewTestRouter(t)
	router.Use(middlewareHTTP.NewTimeoutMiddleware(200 * time.Millisecond))
	router.GET("/fast", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/fast", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
