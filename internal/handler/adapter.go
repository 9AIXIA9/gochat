package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/utils"
	"runtime/debug"
	"time"
)

// Adapter 转换usecase为handler
func Adapter[E domain.ExternalRequest[D], D any](
	usecase domain.Usecase[D],
	cookieSetter func(c *gin.Context, response *domain.Response),
	timeoutDur time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 绑定并校验请求
		hReq := new(E)
		if resp, err := bindParamsAndHandle(c, hReq); err != nil || resp != nil {
			// bindParamsAndHandle 已经记录日志并返回适当响应，直接返回
			return
		}

		// 2. 注入认证信息并转换为领域请求
		injectAuthInfo(c, hReq)
		domainReq := (*hReq).ToDomain()

		// 3. 执行业务逻辑（带超时、recover）
		resp, err := execWithTimeout(c.Request.Context(), timeoutDur, func(ctx context.Context) (*domain.Response, error) {
			return usecase.Execute(ctx, domainReq)
		})

		// 4. 统一处理结果/错误
		if err != nil {
			handleExecError(c, err)
			return
		}
		handleExecSuccess(c, cookieSetter, resp)
	}
}

type panicErr struct {
	Value any
	Stack string
}

func (p panicErr) Error() string {
	return fmt.Sprintf("panic: %v\n%s", p.Value, p.Stack)
}

// bindParamsAndHandle 负责 BindParams 调用及错误处理；
// 若返回非 nil 的 response，表示已向客户端写入响应（比如验证失败或提前返回）
func bindParamsAndHandle[E domain.ExternalRequest[D], D any](c *gin.Context, hReq *E) (*domain.Response, error) {
	resp, err := BindParams(c, hReq)
	if err != nil {
		zap.L().Error("bind params failed", zap.Error(err))
		ResponseError(c)
		return nil, err
	}
	if resp != nil {
		ResponseSuccess(c, resp)
		return resp, nil
	}
	return nil, nil
}

// execWithTimeout 在独立 goroutine 中运行 logic，带 recover 并通过 channel 返回结果或错误。
// 返回 nil, err 表示出错或超时；返回 resp, nil 表示成功。
func execWithTimeout(parentCtx context.Context, timeoutDur time.Duration, logic func(ctx context.Context) (*domain.Response, error)) (*domain.Response, error) {
	ctx, cancel := context.WithTimeout(parentCtx, timeoutDur)
	defer cancel()

	resultCh := make(chan *domain.Response, 1)
	errCh := make(chan error, 1)

	go func() {
		// recover 并把 panic 信息发回主协程
		defer func() {
			if r := recover(); r != nil {
				errCh <- panicErr{
					Value: r,
					Stack: string(debug.Stack()),
				}
			}
		}()

		resp, err := logic(ctx)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- resp
	}()

	select {
	case <-ctx.Done():
		// context 超时或取消
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case resp := <-resultCh:
		return resp, nil
	}
}

// handleExecError 根据错误类型统一处理超时、panic、普通错误
func handleExecError(c *gin.Context, err error) {
	// 1. panicErr
	var p panicErr
	if errors.As(err, &p) {
		zap.L().Error("panic recovered in business logic",
			zap.Any("panic", p.Value),
			zap.String("stack", p.Stack),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		ResponseError(c)
		return
	}

	// 2. 超时或取消
	if utils.IsCanceledOrTimeout(err) {
		ResponseSuccess(c, domain.TimeoutResponse)
		return
	}

	// 3. 其他业务错误
	zap.L().Error("server error", zap.Error(err))
	ResponseError(c)
}

// handleExecSuccess 处理成功响应并设置 cookie（如果有）
func handleExecSuccess(c *gin.Context, cookieSetter func(c *gin.Context, response *domain.Response), resp *domain.Response) {
	if resp != nil {
		if cookieSetter != nil {
			cookieSetter(c, resp)
		}
		ResponseSuccess(c, resp)
		return
	}
	ResponseSuccess(c, domain.DefaultResponse)
}
