package handler_test

import (
	"gochat/internal/delivery/http/handler"
	"gochat/internal/shared/api"
	httptestutil "gochat/pkg/httptest"
	"testing"
)

func TestNewNotFoundHandler(t *testing.T) {
	router := httptestutil.NewTestRouter(t)
	router.NoRoute(handler.NewNotFoundHandler())
	httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
		Path:     "/not-existing-route",
		Method:   "GET",
		Response: api.ResponseNotFound,
	})
}
