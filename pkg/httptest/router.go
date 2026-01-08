package httptest

import (
	"gochat/internal/shared/api"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Path     string
	Method   string
	Body     io.Reader
	Response *api.Response
}

func NewTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func ExecuteRouteTest(t *testing.T, router *gin.Engine, behavior *TestCase) {
	if behavior == nil {
		t.Fatal("behavior is nil")
	}

	req, err := http.NewRequest(behavior.Method, behavior.Path, behavior.Body)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, behavior.Response.String(), w.Body.String())
	assert.Equal(t, behavior.Response.Code.ToHTTPCode(), w.Code)
}
