package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgreeFriendRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required"`
}

func (r *AgreeFriendRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

func NewAgreeFriendRequestHandler(useCase application.AgreeFriendRequestUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *AgreeFriendRequestRequest) *application.AgreeFriendRequestInput {
			return &application.AgreeFriendRequestInput{
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
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "friend request not found"))
			case errors.Is(err, domain.ErrFriendRequestNotForUser):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "friend request not for this user"))
			case errors.Is(err, domain.ErrFriendRequestHasBeenHandled):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "friend request has been handled"))
			default:
				zap.L().Error("SendFriendRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
