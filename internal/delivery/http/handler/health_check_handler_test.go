package handler_test

import (
	"gochat/internal/delivery/http/handler"
	"gochat/internal/shared/api"
	"net/http"
	"net/http/httptest"
	"testing"

	httptestutil "gochat/pkg/httptest"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthCheckHandler(t *testing.T) {
	router := httptestutil.SetupRouter(t)

	router.GET("/healthz", handler.NewHealthCheckHandler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, api.NewResponseWithMessage(api.CodeSuccess, "server is healthy").String(), w.Body.String())
}
