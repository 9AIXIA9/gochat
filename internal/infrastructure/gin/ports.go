package gin

import (
	"context"

	"github.com/gin-gonic/gin"
)

type RequestPointers[Request any] interface {
	*Request
	Bindable
}

type Bindable interface {
	Bind(ginContext *gin.Context) error
}

type Validator interface {
	Validate(ctx context.Context, model any) (string, error)
}
