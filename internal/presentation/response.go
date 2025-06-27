package presentation

import (
	"go.uber.org/zap"
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, data ...interface{}) {
	var responseData interface{}
	if len(data) == 0 {
		responseData = nil
	} else if len(data) == 1 {
		responseData = data[0]
	} else {
		responseData = data
	}

	c.JSON(domain.CodeSuccess.ToHTTP(), domain.Message{
		Code: domain.CodeSuccess,
		Msg:  domain.CodeSuccess.Msg(),
		Data: responseData,
	})
}

func ResponseError(c *gin.Context, code domain.ResCode, msg string, field ...zap.Field) {
	if len(msg) != 0 {
		zap.L().Error(msg, field...)
	}

	c.JSON(code.ToHTTP(), domain.Message{
		Code: code,
		Msg:  code.Msg(),
	})
}
