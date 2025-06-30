package handler

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, message *domain.Message) {
	c.JSON(message.Code.ToHTTP(), message)
}

func ResponseError(c *gin.Context) {
	c.JSON(domain.CodeServerBusy.ToHTTP(), domain.Message{
		Code: domain.CodeServerBusy,
		Msg:  domain.CodeServerBusy.Msg(),
	})
}
