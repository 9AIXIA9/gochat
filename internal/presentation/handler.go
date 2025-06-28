package presentation

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

// HandlerAdapter 将业务逻辑处理函数转换为gin.HandlerFunc
func HandlerAdapter[Req any](logicFn func(*Req) (*domain.Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 绑定请求参数
		var req Req
		if msg, err := BindParams(c, &req); err != nil {
			zap.L().Error("bind param error", zap.Error(err))
			ResponseError(c)
			return
		} else if msg != nil {
			//返回提示信息
			ResponseSuccess(c, msg)
			return
		}

		// 注入认证信息
		if authReq, ok := any(&req).(domain.AuthContext); ok {
			userNumber := utils.GetCurrentUser(c)
			authReq.SetUserNumber(userNumber)
		}

		//todo err 元数据暂未处理
		resp, err := logicFn(&req)
		if err != nil {
			zap.L().Error("logic error", zap.Error(err))
			ResponseError(c)
			return
		}

		if resp != nil {
			//指定成功响应
			ResponseSuccess(c, resp)
			return
		}

		//默认成功响应
		ResponseSuccess(c, domain.NewSuccessResponse())
	}
}
