package middleware_test

import (
	middlewareHTTP "gochat/internal/delivery/http/middleware"
	"gochat/internal/shared/api"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRecoverMiddleware(t *testing.T) {
	t.Parallel()

	router := httptestutil.NewTestRouter(t)
	router.Use(middlewareHTTP.NewRecoverMiddleware())
	router.GET("/panic", func(_ *gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, api.ResponseServerError.String(), w.Body.String())
	assert.Equal(t, api.ResponseServerError.Code.ToHTTPCode(), w.Code)
}
