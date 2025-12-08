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

type RefuseMemberRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *RefuseMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

func NewRefuseMemberRequestHandler(useCase application.RefuseMemberRequestUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *RefuseMemberRequestRequest) *application.RefuseMemberRequestInput {
			return &application.RefuseMemberRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.Response(ginContext, sharedHttp.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "request is not found")
			case errors.Is(err, domain.ErrNotAdmin):
				sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "you have no permission to refuse this request")
			default:
				zap.L().Error("RefuseMemberRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
