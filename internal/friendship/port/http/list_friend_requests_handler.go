package http

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/dto"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

//TODO 所有的limit要重新限制

const defaultFriendRequestsLimit = 20

type ListFriendRequestsRequest struct {
	UserID kernel.UserID      `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	BaseID kernel.OperationID `form:"base_id" example:"019b593b-462e-74d6-bfda-0e103a172191"` // 用于分页游标
	Limit  int                `form:"limit" validate:"omitempty,min=1,max=100" example:"20"`  // 每页条数
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

// NewListFriendRequestsHandler 获取该处理的好友请求列表
// @Summary      获取该处理的好友请求列表
// @Description  获取当前登录用户相关的好友请求列表，可基于 base_id 游标和 limit 分页
// @Tags         Friendship
// @Security     BearerAuth
// @Param        base_id  query     string    false "分页游标，返回该ID之前的记录"    example("019b593b-462e-74d6-bfda-0e103a172190")
// @Param        limit    query     int    false "分页大小，默认20，最大100"   minimum(1) maximum(100) default(20) example(50)
// @Success      200      {object}  api.Response{data=ListFriendRequestsResponseData} "成功返回好友请求列表"
// @Router       /friendship-requests/ [get]
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
	)
}
