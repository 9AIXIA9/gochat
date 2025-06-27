package logger

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinOption() gin.OptionFunc {
	return func(r *gin.Engine) {
		r.Use(GinLogger(), GinRecovery())
	}
}

// GinLogger 自定义接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 用于恢复可能出现的panic
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rawErr := recover(); rawErr != nil {
				// 统一将 panic 值转为 error 类型
				var err error
				switch v := rawErr.(type) {
				case string:
					err = errors.New(v)
				case error:
					err = v
				default:
					err = fmt.Errorf("unknown panic: %v", v)
				}

				// 处理网络连接错误
				var ne *net.OpError
				if errors.As(err, &ne) {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if isConnectionError(se) {
							logBrokenPipe(c, err)
							c.Abort()
							return
						}
					}
				}

				// 记录其他错误详情
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := make(map[string][]string)
				for k, v := range c.Request.Header {
					headers[k] = v
				}

				zap.L().Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Error(err), // 使用安全类型记录
					zap.String("request", string(httpRequest)),
					zap.Any("headers", headers),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// 发生 "broken pipe" 或 "connection reset by peer"
// 这类网络错误时，logBrokenPipe 会被调用。
// 这些错误通常发生在客户端突然断开连接的情况下
func logBrokenPipe(c *gin.Context, err interface{}) {
	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	zap.L().Error(c.Request.URL.Path,
		zap.Any("error", err),
		zap.String("request", string(httpRequest)),
	)
}

// 识别特定类型的网络连接错误
// "broken pipe"（管道破裂）
// "connection reset by peer"（对方重置连接）
func isConnectionError(se *os.SyscallError) bool {
	errStr := se.Error()
	//TOLower将字母转化为小写
	return strings.Contains(strings.ToLower(errStr), "broken pipe") ||
		strings.Contains(strings.ToLower(errStr), "connection reset by peer")
}
