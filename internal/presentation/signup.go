package presentation

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

type SignupHandler struct {
	signupUsecase domain.SignupUsecase
	logger        domain.Logger
}

func SignupHandlerFunc(signupUsecase domain.SignupUsecase, logger domain.Logger) gin.HandlerFunc {
	return (&SignupHandler{
		signupUsecase: signupUsecase,
		logger:        logger,
	}).Signup
}

func (h *SignupHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("注册参数绑定失败: %v", err)
		ResponseError(c, domain.CodeInvalidParam, err.Error())
		return
	}

	// 加密密码
	hashedPassword := h.signupUsecase.EncryptPwd(req.Password)

	// 生成用户号码
	userNumber := h.signupUsecase.GenerateNumber()

	// 创建用户
	user := domain.User{
		Number:  userNumber,
		Name:    req.Name,
		PwdHash: hashedPassword,
	}

	if err := h.signupUsecase.CreateUser(user); err != nil {
		h.logger.Error("创建用户失败: %v", err)
		ResponseError(c, domain.CodeUserExist, domain.CodeUserExist.Msg())
		return
	}

	// 返回用户信息
	ResponseSuccess(c, gin.H{
		"user_number": userNumber,
		"name":        req.Name,
	})
}
