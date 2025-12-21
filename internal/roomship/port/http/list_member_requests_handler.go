package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/dto"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
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

// NewListMemberRequestsHandler 获取房间成员请求列表
// @Summary      获取房间成员请求列表
// @Description  获取当前登录用户相关的房间成员请求列表，可基于 base_id 游标和 limit 分页
// @Tags         Roomship
// @Security     BearerAuth
// @Produce      json
// @Param        base_id  query       kernel.OperationID     false "分页游标，返回该ID之前的记录"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  api.Response{data=ListMemberRequestsResponseData} "成功返回成员请求列表"
// @Failure      400      {object}  api.Response         "请求参数错误"
// @Failure      401      {object}  api.Response         "未认证"
// @Failure      500      {object}  api.Response         "服务器内部错误"
// @Router       /roomship/request [get]
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
			ginutils.ResponseSuccessWithData(ginContext, &ListMemberRequestsResponseData{
				Requests: dto.ToMemberRequestDTOs(output.Requests),
			})
		},
	)
}
