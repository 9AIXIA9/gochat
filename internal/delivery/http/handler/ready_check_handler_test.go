package handler_test

import (
	"gochat/internal/delivery/http/handler"
	"gochat/internal/shared/api"
	httptestutil "gochat/pkg/httptest"
	"testing"
)

func TestNewReadyCheckHandler_Success(t *testing.T) {
	router := httptestutil.NewTestRouter(t)
	router.GET("/readyz", handler.NewReadyCheckHandler(func() bool {
		return true
	}))
	httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
		Path:     "/readyz",
		Method:   "GET",
		Response: api.NewResponseWithMessage(api.CodeSuccess, "server is ready"),
	})
}

func TestNewReadyCheckHandler_NotReady(t *testing.T) {
	router := httptestutil.NewTestRouter(t)
	router.GET("/readyz", handler.NewReadyCheckHandler(func() bool {
		return false
	}))
	httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
		Path:     "/readyz",
		Method:   "GET",
		Response: api.NewResponseWithMessage(api.CodeServiceUnavailable, "server is not ready"),
	})
}
