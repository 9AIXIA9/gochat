package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application/usecase"
	"gochat/internal/social/domain"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateRoomRequest struct {
	Owner          kernel.UserID   `json:"-" validate:"required"`
	MaxMemberCount int             `json:"max_member_count" validate:"required,min=2,max=100"`
	Password       domain.Password `json:"password" validate:"max=100"`
}

type CreateRoomResponseData struct {
	RoomNumber domain.RoomNumber `json:"room_number"`
}

func (r *CreateRoomRequest) Bind(ginContext *gin.Context) error {
	owner := ginutils.GetUserID(ginContext)
	r.Owner = owner
	return ginContext.ShouldBind(r)
}

func NewCreateRoomHandler(useCase usecase.CreateRoomUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *CreateRoomRequest) *usecase.CreateRoomInput {
			return &usecase.CreateRoomInput{
				Owner:          request.Owner,
				MaxMemberCount: request.MaxMemberCount,
				Password:       request.Password,
			}
		},
		func(ginContext *gin.Context, output *usecase.CreateRoomOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&CreateRoomResponseData{
				RoomNumber: output.RoomNumber,
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "recipient not found"))
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
