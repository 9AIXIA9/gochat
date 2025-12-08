package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
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
	UserID         kernel.UserID   `json:"-" validate:"required"`
	MaxMemberCount int             `json:"max_member_count" validate:"gte=0,lte=200"`
	Password       domain.Password `json:"password" validate:"max=100"`
}

func (r *CreateRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

func NewCreateRoomHandler(useCase application.CreateRoomUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *CreateRoomRequest) *application.CreateRoomInput {
			return &application.CreateRoomInput{
				UserID:         request.UserID,
				MaxMemberCount: request.MaxMemberCount,
				Password:       request.Password,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("CreateRoomHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
