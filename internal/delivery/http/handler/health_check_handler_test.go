package handler_test

import (
	"gochat/internal/delivery/http/handler"
	"gochat/internal/shared/api"
	"testing"

	httptestutil "gochat/pkg/httptest"
)

func TestNewHealthCheckHandler(t *testing.T) {
	router := httptestutil.NewTestRouter(t)
	router.GET("/healthz", handler.NewHealthCheckHandler())
	httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
		Path:     "/healthz",
		Method:   "GET",
		Response: api.NewResponseWithMessage(api.CodeSuccess, "server is healthy"),
	})
}
