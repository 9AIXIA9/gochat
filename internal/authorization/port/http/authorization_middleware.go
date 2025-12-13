package http

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/http"
	"gochat/pkg/utils"

	"github.com/gin-gonic/gin"

	"strings"
)

func NewAuthorizationMiddleware(useCase application.ParseAccessTokenUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// 从请求头获取token
		var accessToken domain.AccessToken
		if authorizationHeader := ginContext.Request.Header.Get("Authorization"); authorizationHeader != "" {
			// Bearer token格式
			parts := strings.SplitN(authorizationHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				ginutils.Response(ginContext, http.CodeInvalidToken)
				ginContext.Abort()
				return
			}
			accessToken = domain.AccessToken(parts[1])
		} else if authorizationQuery := ginContext.Query("Authorization"); authorizationQuery != "" {
			accessToken = domain.AccessToken(authorizationQuery)
		} else {
			ginutils.Response(ginContext, http.CodeInvalidToken)
			ginContext.Abort()
			return
		}

		//解析 token
		output, err := useCase.Execute(ginContext.Request.Context(), &application.ParseAccessTokenInput{AccessToken: accessToken})
		if err != nil {
			ginutils.Response(ginContext, http.CodeInvalidToken)
			ginContext.Abort()
			return
		}

		// 写入到 gin.Context 供 gin handlers 使用
		ginutils.SetUserID(ginContext, output.UserID)
		// 同步写入到 request.Context，供标准 http.Handler 使用
		ginContext.Request = ginContext.Request.WithContext(utils.SetUserID(ginContext.Request.Context(), output.UserID))
		ginContext.Next()
	}
}
