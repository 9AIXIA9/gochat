package presentation

import (
	"go.uber.org/zap"
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

type SignupHandler struct {
	signupUsecase domain.SignupUsecase
}

func SignupHandlerFunc(signupUsecase domain.SignupUsecase) gin.HandlerFunc {
	return (&SignupHandler{
		signupUsecase: signupUsecase,
	}).Signup
}

func (h *SignupHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, domain.CodeInvalidParam, "")
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
		ResponseError(c, domain.CodeUserExist, "create user failed", zap.Error(err))
		return
	}

	// 返回用户信息
	ResponseSuccess(c, gin.H{
		"user_number": userNumber,
		"name":        req.Name,
	})
}
