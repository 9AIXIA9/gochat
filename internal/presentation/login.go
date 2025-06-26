package presentation

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUsecase domain.LoginUsecase
	logger       domain.Logger
}

func LoginHandlerFunc(loginUsecase domain.LoginUsecase, logger domain.Logger) gin.HandlerFunc {
	return (&LoginHandler{
		loginUsecase: loginUsecase,
		logger:       logger,
	}).Login
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("登录参数绑定失败: %v", err)
		ResponseError(c, domain.CodeInvalidParam, err.Error())
		return
	}

	// 验证用户名密码
	if err := h.loginUsecase.CheckPwd(req.Number, req.Password); err != nil {
		h.logger.Error("密码验证失败: %v", err)
		ResponseError(c, domain.CodeWrongPassword, domain.CodeWrongPassword.Msg())
		return
	}

	// 生成token
	token := h.loginUsecase.GenerateToken()
	if token == "" {
		h.logger.Error("生成token失败")
		ResponseError(c, domain.CodeServerBusy, domain.CodeServerBusy.Msg())
		return
	}

	// 返回token
	ResponseSuccess(c, gin.H{
		"token":       token,
		"user_number": req.Number,
	})
}
