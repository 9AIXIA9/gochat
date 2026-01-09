//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package gin

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Bindable interface {
	Bind(ginContext *gin.Context) error
}

type Validator interface {
	Validate(ctx context.Context, model any) (string, error)
}
