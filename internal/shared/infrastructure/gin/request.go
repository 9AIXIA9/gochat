package gin

import (
	"github.com/gin-gonic/gin"
)

type Bindable interface {
	Bind(ginContext *gin.Context) error
}
