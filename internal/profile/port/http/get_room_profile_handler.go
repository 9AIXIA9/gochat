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

type GetRoomProfileRequest struct {
	RoomID kernel.RoomID `uri:"room_id" validate:"required"`
}

func (r *GetRoomProfileRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type GetRoomProfileResponseData struct {
	Profile *dto.RoomProfile `json:"profile"`
}

func NewGetRoomProfileHandler(useCase application.GetRoomProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetRoomProfileRequest) *application.GetRoomProfileInput {
			return &application.GetRoomProfileInput{
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, output *application.GetRoomProfileOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&GetRoomProfileResponseData{
				Profile: dto.ToRoomProfileDTO(output.Profile),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeNotFound, "room is not found"))
			default:
				zap.L().Error("GetRoomProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
