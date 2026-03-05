package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	"gochat/internal/infrastructure/breaker"
	"gochat/internal/shared/api"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreakMiddleware(t *testing.T) {
	t.Parallel()

	router := httptestutil.NewTestRouter(t)
	router.Use(middlewareHTTP.NewCircuitBreakMiddleware(&breaker.Config{
		Name:           "test-breaker",
		MaxRequests:    1,
		Interval:       time.Second,
		BucketPeriod:   time.Second,
		Timeout:        time.Second,
		ErrorRate:      0.5,
		ErrorThreshold: 1,
	}))
	router.GET("/fail", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	router.ServeHTTP(w, req)

	expect := api.NewResponseWithMessage(api.CodeServiceUnavailable, "service unavailable")
	assert.Equal(t, expect.String(), w.Body.String())
	assert.Equal(t, expect.Code.ToHTTPCode(), w.Code)
}
