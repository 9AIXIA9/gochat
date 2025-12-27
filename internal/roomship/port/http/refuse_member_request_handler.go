package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
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

// NewRefuseMemberRequestHandler 拒绝入群请求
// @Summary      拒绝入群请求
// @Description  拒绝指定房间成员请求
// @Tags         Roomship
// @Security     BearerAuth
// @Param        request_id  path      kernel.OperationID                   true  "成员请求ID"
// @Success      200         {object}  api.Response "拒绝成功"
// @Router       /rooms/requests/{request_id}/refuse [put]
func NewRefuseMemberRequestHandler(useCase application.RefuseMemberRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
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
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
