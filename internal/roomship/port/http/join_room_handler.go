package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JoinRoomRequest struct {
	UserID     kernel.UserID     `json:"-" validate:"required"`
	RoomNumber kernel.RoomNumber `json:"room_number" validate:"required"`
	Password   domain.Password   `json:"password" validate:"max=100"`
}

func (r *JoinRoomRequest) Bind(ginContext *gin.Context) error {
	id := ginutils.GetUserID(ginContext)
	r.UserID = id
	return ginContext.ShouldBind(r)
}

func NewJoinRoomHandler(useCase application.JoinRoomUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *JoinRoomRequest) *application.JoinRoomInput {
			return &application.JoinRoomInput{
				UserID:     request.UserID,
				RoomNumber: request.RoomNumber,
				Password:   request.Password,
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
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "room not found"))
			case errors.Is(err, myErrors.ErrExceedMaxValue):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeMaxReached, "room members exceed max value"))
			case errors.Is(err, myErrors.ErrInvalidCredential):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "invalid room password"))
			default:
				zap.L().Error("join room handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
