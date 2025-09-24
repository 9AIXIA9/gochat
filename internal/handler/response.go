package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
)

func ResponseSuccess(c *gin.Context, response *domain.Response) {
	c.JSON(response.Code.ToHTTP(), response)
}

func ResponseError(c *gin.Context) {
	c.JSON(domain.CodeServerBusy.ToHTTP(), domain.Response{
		Code: domain.CodeServerBusy,
		Msg:  domain.CodeServerBusy.Msg(),
	})
}
