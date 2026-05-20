package http

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	"gochat/pkg/ctxutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuthorizationMiddleware(useCase application.ParseAccessTokenUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		accessToken, ok := accessTokenFromRequest(ginContext)
		if !ok {
			ginutils.Response(ginContext, api.CodeUnauthorized)
			ginContext.Abort()
			return
		}

		//解析 token
		output, err := useCase.Execute(ginContext.Request.Context(), &application.ParseAccessTokenInput{AccessToken: accessToken})
		if err != nil {
			ginutils.Response(ginContext, api.CodeUnauthorized)
			ginContext.Abort()
			return
		}

		// 写入到 request.Context，供标准 http.Handler 使用
		ginContext.Request = ginContext.Request.WithContext(ctxutil.WithUserID(ginContext.Request.Context(), output.UserID))
		ginContext.Next()
	}
}

func accessTokenFromRequest(ginContext *gin.Context) (domain.AccessToken, bool) {
	if token := strings.TrimSpace(ginContext.Query("access_token")); token != "" {
		return domain.AccessToken(token), true
	}

	authorizationHeader := strings.TrimSpace(ginContext.Request.Header.Get("Authorization"))
	if authorizationHeader == "" {
		return "", false
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return domain.AccessToken(token), true
}
