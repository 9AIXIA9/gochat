package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/http"
	ginutils "gochat/internal/shared/infrastructure/gin"
	"gochat/internal/shared/kernel"

	"strings"
)

const UserIDKey = "user_id"

func Authorization(useCase kernel.ParseAccessTokenUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// 从请求头获取token
		var accessToken kernel.AccessToken
		if authorizationHeader := ginContext.Request.Header.Get("Authorization"); authorizationHeader != "" {
			// Bearer token格式
			parts := strings.SplitN(authorizationHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				ginutils.Response(ginContext, http.ResponseInvalidToken)
				ginContext.Abort()
				return
			}
			accessToken = kernel.AccessToken(parts[1])
		} else if authorizationQuery := ginContext.Query("Authorization"); authorizationQuery != "" {
			accessToken = kernel.AccessToken(authorizationQuery)
		} else {
			ginutils.Response(ginContext, http.ResponseInvalidToken)
			ginContext.Abort()
			return
		}

		// 解析token
		if output, err := useCase.Execute(ginContext.Request.Context(), &kernel.ParseAccessTokenInput{AccessToken: accessToken}); err != nil {
			if errors.Is(err, myErrors.ErrInvalidCredential) {
				ginutils.Response(ginContext, http.ResponseInvalidToken)
				ginContext.Abort()
				return
			} else {
				ginutils.Response(ginContext, http.ResponseInvalidToken)
				ginContext.Abort()
				return
			}
		} else {
			ginContext.Set(UserIDKey, output.UserID)
			ginContext.Next()
		}
	}
}
