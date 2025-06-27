package presentation

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUsecase domain.LoginUsecase
}

func LoginHandlerFunc(loginUsecase domain.LoginUsecase) gin.HandlerFunc {
	return (&LoginHandler{
		loginUsecase: loginUsecase,
	}).Login
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, domain.CodeInvalidParam, "")
		return
	}

	// 验证用户名密码
	if err := h.loginUsecase.CheckPwd(req.Number, req.Password); err != nil {
		ResponseError(c, domain.CodeWrongPassword, "wrong password or userNumber")
		return
	}

	// 生成token
	token := h.loginUsecase.GenerateToken()
	if token == "" {
		ResponseError(c, domain.CodeServerBusy, "")
		return
	}

	// 返回token
	ResponseSuccess(c, gin.H{
		"token":       token,
		"user_number": req.Number,
	})
}
