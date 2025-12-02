package http

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/dto"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultLimit = 20

type ListFriendRequestsRequest struct {
	UserID kernel.UserID      `json:"-" validate:"required"`
	BaseID kernel.OperationID `form:"base_id"`                                  // 用于分页游标
	Limit  int                `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListFriendRequestsRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultLimit
	}
	return nil
}

type ListFriendRequestsResponseData struct {
	Requests []*dto.FriendRequest `json:"requests,omitempty"`
}

func NewListFriendRequestsHandler(useCase application.ListFriendRequestsUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListFriendRequestsRequest) *application.ListFriendRequestsInput {
			return &application.ListFriendRequestsInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListFriendRequestsOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListFriendRequestsResponseData{
				Requests: dto.ToFriendRequestDTOs(output.Requests),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch err.(type) {
			//TODO handle specific errors
			default:
				zap.L().Error("SendFriendRequestHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
