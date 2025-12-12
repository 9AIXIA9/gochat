package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateUserProfileRequest struct {
	UserID      kernel.UserID      `json:"-" validate:"required"`
	Name        string             `json:"name"`
	Gender      kernel.Gender      `json:"gender"`
	Email       kernel.Email       `json:"email"`
	PhoneNumber kernel.PhoneNumber `json:"phone_number"`
	Address     kernel.Address     `json:"address"`
	Sign        string             `json:"sign"`
}

func (r *UpdateUserProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

func NewUpdateUserProfileHandler(useCase application.UpdateUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *UpdateUserProfileRequest) *application.UpdateUserProfileInput {
			return &application.UpdateUserProfileInput{
				UserID:      request.UserID,
				Name:        request.Name,
				Gender:      request.Gender,
				Email:       request.Email,
				PhoneNumber: request.PhoneNumber,
				Address:     request.Address,
				Sign:        request.Sign,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeNotFound, "user profile not found"))
			default:
				zap.L().Error("UpdateUserProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
