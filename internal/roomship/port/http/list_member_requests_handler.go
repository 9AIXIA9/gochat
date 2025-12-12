package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/dto"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultMemberRequestsLimit = 20

type ListMemberRequestsRequest struct {
	UserID kernel.UserID      `json:"-" validate:"required"`
	BaseID kernel.OperationID `form:"base_id"`                                  // 用于分页游标
	Limit  int                `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListMemberRequestsRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultMemberRequestsLimit
	}
	return nil
}

type ListMemberRequestsResponseData struct {
	Requests []*dto.MemberRequest `json:"requests,omitempty"`
}

func NewListMemberRequestsHandler(useCase application.ListMemberRequestsUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListMemberRequestsRequest) *application.ListMemberRequestsInput {
			return &application.ListMemberRequestsInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListMemberRequestsOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListMemberRequestsResponseData{
				Requests: dto.ToMemberRequestDTOs(output.Requests),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			default:
				zap.L().Error("ListMemberRequestsHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
