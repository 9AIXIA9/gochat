package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
)

// Adapter 将HTTP请求处理转换为业务逻辑处理函数
func Adapter[E domain.ExternalRequest[D], D any](usecase domain.Usecase[D]) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求对象
		hReq := new(E)

		// 绑定请求参数
		if resp, err := BindParams(c, hReq); err != nil {
			zap.L().Error("bind params failed", zap.Error(err))
			ResponseError(c)
			return
		} else if resp != nil {
			ResponseSuccess(c, resp)
			return
		}

		// 注入认证信息
		injectAuthInfo(c, hReq)

		// 转换为领域请求并执行逻辑
		domainReq := (*hReq).ToDomain()
		executeLogic(c, func() (*domain.Response, error) {
			return usecase.Logic(domainReq)
		})
	}
}

func executeLogic(c *gin.Context, logic func() (*domain.Response, error)) {
	resp, err := logic()
	if err != nil {
		zap.L().Error("server error", zap.Error(err))
		ResponseError(c)
		return
	}

	// 处理响应
	if resp != nil {
		ResponseSuccess(c, resp)
	} else {
		ResponseSuccess(c, domain.NewSuccessResponse())
	}
}
