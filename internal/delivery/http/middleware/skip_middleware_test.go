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

func TestSkipMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		path             string
		route            string
		skippedPaths     []string
		expectMiddleware bool
	}{
		{
			name:             "skip exact path",
			path:             "/healthz",
			route:            "/healthz",
			skippedPaths:     []string{"/healthz"},
			expectMiddleware: false,
		},
		{
			name:             "skip wildcard path",
			path:             "/swagger/index.html",
			route:            "/swagger/*any",
			skippedPaths:     []string{"/swagger/*any"},
			expectMiddleware: false,
		},
		{
			name:             "run middleware for business path",
			path:             "/api/v1/ping",
			route:            "/api/v1/ping",
			skippedPaths:     []string{"/healthz"},
			expectMiddleware: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := httptestutil.NewTestRouter(t)

			called := false
			mw := middlewareHTTP.SkipMiddleware(tt.skippedPaths, func(c *gin.Context) {
				called = true
				c.Next()
			})

			router.Use(mw)
			router.GET(tt.route, func(c *gin.Context) { c.Status(http.StatusOK) })

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expectMiddleware, called)
		})
	}
}
