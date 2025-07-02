package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
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
