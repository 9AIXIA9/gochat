package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type AgreeMemberRequestRequest struct {
	UserID    kernel.UserID      `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RequestID kernel.OperationID `uri:"request_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172192"`
}

func (r *AgreeMemberRequestRequest) Bind(ginContext *gin.Context) error {
	userID := ctxutil.UserIDFrom(ginContext.Request.Context())
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewAgreeMemberRequestHandler 同意入群请求
// @Summary      同意入群请求
// @Description  同意指定房间成员请求，将对方加入房间
// @Tags         Roomship
// @Security     BearerAuth
// @Param        request_id  path      string    true  "成员请求ID"    example(019b593b-462e-74d6-bfda-0e103a172190)
// @Success      200         {object}  api.Response "同意成功"
// @Router       /rooms/requests/{request_id}/agree [put]
func NewAgreeMemberRequestHandler(useCase application.AgreeMemberRequestUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *AgreeMemberRequestRequest) *application.AgreeMemberRequestInput {
			return &application.AgreeMemberRequestInput{
				UserID:    request.UserID,
				RequestID: request.RequestID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
