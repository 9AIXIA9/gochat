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

type CreateRoomRequest struct {
	OwnerID        kernel.UserID   `json:"-" validate:"required"`
	MaxMemberCount int             `json:"max_member_count" validate:"min=2,max=100"`
	Password       domain.Password `json:"password" validate:"max=100"`
}

type CreateRoomResponseData struct {
	RoomNumber kernel.RoomNumber `json:"room_number"`
}

func (r *CreateRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.OwnerID = userID
	return ginContext.ShouldBind(r)
}

func NewCreateRoomHandler(useCase application.CreateRoomUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *CreateRoomRequest) *application.CreateRoomInput {
			return &application.CreateRoomInput{
				OwnerID:        request.OwnerID,
				MaxMemberCount: request.MaxMemberCount,
				Password:       request.Password,
			}
		},
		func(ginContext *gin.Context, output *application.CreateRoomOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&CreateRoomResponseData{
				RoomNumber: output.RoomNumber,
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrInvalidNumber):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "max member number is invalid"))
			default:
				zap.L().Error("create room handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
