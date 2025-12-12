package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserProfileRequest struct {
	UserID kernel.UserID `uri:"user_id" validate:"required"`
}

func (r *GetUserProfileRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type GetUserProfileResponseData struct {
	Profile *dto.UserProfile `json:"profile"`
}

func NewGetUserProfileHandler(useCase application.GetUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetUserProfileRequest) *application.GetUserProfileInput {
			return &application.GetUserProfileInput{
				UserID: request.UserID,
			}
		},
		func(ginContext *gin.Context, output *application.GetUserProfileOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&GetUserProfileResponseData{
				Profile: dto.ToUserProfileDTO(output.Profile),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeNotFound, "user is not found"))
			default:
				zap.L().Error("GetUserProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
