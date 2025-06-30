package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/types"
	"gochat/internal/utils/timeout"
)

// Adapter 将HTTP请求处理转换为业务逻辑处理函数
func Adapter[E domain.ExternalRequest[D], D any](usecase domain.Usecase[D]) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求对象
		hReq := new(E)

		// 绑定请求参数
		resp, err := BindParams(c, hReq)
		if err != nil {
			zap.L().Error("bind params failed", zap.Error(err))
			ResponseError(c)
			return
		}

		if resp != nil {
			ResponseSuccess(c, resp)
			return
		}

		// 注入认证信息
		injectAuthInfo(c, hReq)

		// 转换为领域请求并执行逻辑
		domainReq := (*hReq).ToDomain()
		executeFn(c, func() (*domain.Message, error) {
			return usecase.Logic(c.Request.Context(), domainReq)
		})
	}
}

func executeFn(c *gin.Context, logic func() (*domain.Message, error)) {
	resp, err := logic()
	if err != nil {
		if timeout.IsCanceledOrTimeout(err) {
			ResponseSuccess(c, types.TimeoutResponse)
			return
		}

		zap.L().Error("server error", zap.Error(err))
		ResponseError(c)
		return
	}

	// 处理响应
	if resp != nil {
		ResponseSuccess(c, resp)
	} else {
		ResponseSuccess(c, types.DefaultResponse)
	}
}
