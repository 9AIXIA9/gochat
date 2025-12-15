package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/dto"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultFriendRequestsLimit = 20

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
		r.Limit = defaultFriendRequestsLimit
	}
	return nil
}

type ListFriendRequestsResponseData struct {
	Requests []*dto.FriendRequest `json:"requests,omitempty"`
}

// NewListFriendRequestsHandler 获取好友请求列表
// @Summary      获取好友请求列表
// @Description  获取当前登录用户相关的好友请求列表，可基于 base_id 游标和 limit 分页
// @Tags         Friendship
// @Security     BearerAuth
// @Produce      json
// @Param        base_id  query     int    false "分页游标，返回该ID之前的记录"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  ListFriendRequestsResponseData "成功返回好友请求列表"
// @Failure      400      {object}  api.Response       "请求参数错误"
// @Failure      401      {object}  api.Response       "未认证"
// @Failure      500      {object}  api.Response       "服务器内部错误"
// @Router       /friendship/request [get]
func NewListFriendRequestsHandler(useCase application.ListFriendRequestsUseCase, validator ginutils.Validator) gin.HandlerFunc {
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
			ginutils.ResponseSuccessWithData(ginContext, &ListFriendRequestsResponseData{
				Requests: dto.ToFriendRequestDTOs(output.Requests),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("ListFriendRequestsHandler error", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
