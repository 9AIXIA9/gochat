package presentation

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, message *domain.Response) {
	c.JSON(message.Code.ToHTTP(), message)
}

func ResponseError(c *gin.Context) {
	c.JSON(domain.CodeServerBusy.ToHTTP(), domain.Response{
		Code: domain.CodeServerBusy,
		Msg:  domain.CodeServerBusy.Msg(),
	})
}
