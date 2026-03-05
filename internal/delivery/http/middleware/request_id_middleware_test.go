package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	"gochat/pkg/ctxutil"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequestIDMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		headerRequestID string
		expectSame      bool
	}{
		{
			name:            "generate request id when missing",
			headerRequestID: "",
			expectSame:      false,
		},
		{
			name:            "reuse incoming request id",
			headerRequestID: "req-123",
			expectSame:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := httptestutil.NewTestRouter(t)
			router.Use(middlewareHTTP.NewRequestIDMiddleware())
			router.GET("/test", func(c *gin.Context) {
				rid := ctxutil.RequestIDFrom(c.Request.Context())
				assert.NotEmpty(t, rid)
				c.String(http.StatusOK, rid.String())
			})

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			require.NoError(t, err)
			if tt.headerRequestID != "" {
				req.Header.Set("X-Request-ID", tt.headerRequestID)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			headerRID := w.Header().Get("X-Request-ID")
			require.NotEmpty(t, headerRID)
			assert.Equal(t, headerRID, w.Body.String())
			if tt.expectSame {
				assert.Equal(t, tt.headerRequestID, headerRID)
			}
		})
	}
}
