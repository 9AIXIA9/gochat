package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	"gochat/internal/shared/api"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewTimeoutMiddleware(t *testing.T) {
	t.Parallel()

	router := httptestutil.NewTestRouter(t)
	router.Use(middlewareHTTP.NewTimeoutMiddleware(10 * time.Millisecond))
	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(50 * time.Millisecond)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, api.ResponseTimeout.String(), w.Body.String())
	assert.Equal(t, api.ResponseTimeout.Code.ToHTTPCode(), w.Code)
}
